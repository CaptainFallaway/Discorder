package discord

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// buildSuperProperties constructs the value for the X-Super-Properties header.
func buildSuperProperties(userAgent string) (string, error) {
	payload := map[string]any{
		"os":                          defaultOS,
		"browser":                     defaultBrowser,
		"release_channel":             defaultReleaseChannel,
		"client_version":              defaultClientVersion,
		"os_version":                  defaultOSVersion,
		"os_arch":                     defaultOSArch,
		"app_arch":                    defaultAppArch,
		"system_locale":               defaultLocale,
		"has_client_mods":             false,
		"client_build_number":         defaultBuildNumber,
		"native_build_number":         defaultNativeBuild,
		"browser_user_agent":          userAgent,
		"browser_version":             defaultBrowserVer,
		"os_sdk_version":              defaultOsSdkVersion,
		"client_event_source":         nil,
		"launch_signature":            "",
		"client_heartbeat_session_id": "",
		"client_app_state":            "focused",
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	buf := &bytes.Buffer{}
	encoder := base64.NewEncoder(base64.StdEncoding, buf)
	if _, err := encoder.Write(b); err != nil {
		return "", err
	}
	if err := encoder.Close(); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// buildRequest creates a new HTTP request with the necessary headers
// and returns it. This is used internally to avoid code duplication.
func (dc *DiscordClient) buildRequest(ctx context.Context, method, path string, queries url.Values, body io.ReadCloser) (*http.Request, error) {
	u := &url.URL{Scheme: "https", Host: "discord.com", Path: fmt.Sprintf("/api/%s%s", ApiVersion, path)}
	if queries != nil {
		u.RawQuery = queries.Encode()
	}

	// Build request with context
	req, err := http.NewRequestWithContext(ctx, method, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	// Set headers to match a discord client request
	req.Header = http.Header{
		"Authorization":      []string{dc.token},
		"Accept":             []string{"*/*"},
		"Accept-Encoding":    []string{defaultAcceptEncoding},
		"Accept-Language":    []string{defaultAcceptLanguage},
		"User-Agent":         []string{defaultUserAgent},
		"X-Discord-Locale":   []string{defaultLocale},
		"X-Discord-Timezone": []string{detectTimezone()},
		"Sec-Fetch-Dest":     []string{"empty"},
		"Sec-Fetch-Mode":     []string{"cors"},
		"Sec-Fetch-Site":     []string{"same-origin"},
		"sec-ch-ua":          []string{defaultSecCHUA},
		"sec-ch-ua-mobile":   []string{defaultSecCHUAMobile},
		"sec-ch-ua-platform": []string{defaultSecCHUAPlatform},
		"x-debug-options":    []string{"bugReporterEnabled"},
		"priority":           []string{"u=1, i"},
	}

	// Contextual referer based on path
	req.Header.Set("Referer", refererForPath(path))

	if sp, err := buildSuperProperties(defaultUserAgent); err == nil {
		req.Header.Set("X-Super-Properties", sp)
	}

	if body != nil {
		req.Body = body
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

// detectTimezone returns an IANA-like timezone string if possible, defaulting to UTC.
func detectTimezone() string {
	if tz := os.Getenv("TZ"); tz != "" {
		return tz
	}
	loc := time.Now().Location()
	if loc != nil && loc.String() != "Local" {
		return loc.String()
	}
	return "UTC"
}

// refererForPath returns a plausible Referer URL for a given API path.
func refererForPath(path string) string {
	referer := "https://discord.com/channels/@me"
	if strings.HasPrefix(path, "/guilds/") {
		parts := strings.Split(path, "/")
		if len(parts) >= 3 && parts[2] != "" {
			return "https://discord.com/channels/" + parts[2]
		}
		return referer
	}
	if strings.HasPrefix(path, "/channels/") {
		parts := strings.Split(path, "/")
		if len(parts) >= 3 && parts[2] != "" {
			return "https://discord.com/channels/@me/" + parts[2]
		}
	}
	return referer
}

// decodeBody wraps the response body with the appropriate decompressor based on Content-Encoding.
func decodeBody(resp *http.Response) (io.ReadCloser, error) {
	enc := strings.ToLower(strings.TrimSpace(resp.Header.Get("Content-Encoding")))
	switch enc {
	case "gzip":
		r, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}
		return r, nil
	case "deflate":
		r, err := zlib.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}
		return r, nil
	default:
		return resp.Body, nil
	}
}
