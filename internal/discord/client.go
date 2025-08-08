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

	"golang.org/x/term"
)

const (
	ApiVersion   = "v9"
	MessageLimit = "100" // Max is 100
)

// Defaults to resemble the Discord desktop client on Windows per provided headers
const (
	defaultUserAgent       = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) discord/1.0.9202 Chrome/134.0.6998.205 Electron/35.3.0 Safari/537.36"
	defaultLocale          = "en-US"
	defaultAcceptLanguage  = "en-US,en-SE;q=0.9,en-GB;q=0.8"
	defaultAcceptEncoding  = "gzip, deflate"
	defaultSecCHUA         = `"Not:A-Brand";v="24", "Chromium";v="134"`
	defaultSecCHUAMobile   = "?0"
	defaultSecCHUAPlatform = `"Windows"`

	// X-Super-Properties related
	defaultBrowser        = "Discord Client"
	defaultBrowserVer     = "35.3.0" // Electron version
	defaultOS             = "Windows"
	defaultOSVersion      = "10.0.26100"
	defaultOSArch         = "x64"
	defaultAppArch        = "x64"
	defaultReleaseChannel = "stable"
	defaultClientVersion  = "1.0.9202"
	defaultBuildNumber    = 429117
	defaultNativeBuild    = 66976
	defaultOsSdkVersion   = "26100"
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

// DiscordClient represents a client for interacting with the Discord API
type DiscordClient struct {
	token  string
	client *http.Client
}

func NewDiscordClient(token string) *DiscordClient {
	return &DiscordClient{token: token, client: &http.Client{}}
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

	// Set headers to match provided set
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

// Simple request with just method and path (context-aware)
func (dc *DiscordClient) Request(ctx context.Context, method, path string) (io.ReadCloser, error) {
	return dc.RequestWithOptions(ctx, method, path, nil, nil)
}

// Full request with queries and body (context-aware)
func (dc *DiscordClient) RequestWithOptions(ctx context.Context, method, path string, queries url.Values, body io.ReadCloser) (io.ReadCloser, error) {
	req, err := dc.buildRequest(ctx, method, path, queries, body)
	if err != nil {
		return nil, err
	}

	// Only print debug info if output is going to a terminal (not piped)
	if term.IsTerminal(int(os.Stdout.Fd())) {
		fmt.Println("Making request:", req.Method, req.URL.String())
	}

	resp, err := dc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Read response body to include in error (limited)
		defer resp.Body.Close()
		limited := io.LimitReader(resp.Body, 8192)
		b, _ := io.ReadAll(limited)
		return nil, fmt.Errorf("request failed: %s: %s", resp.Status, strings.TrimSpace(string(b)))
	}

	// Decode compressed responses if any
	decoded, err := decodeBody(resp)
	if err != nil {
		// Fallback to raw body on decode error
		return resp.Body, nil
	}
	return decoded, nil
}

// GetAllRelationships retrieves all relationships for the authenticated user
func (dc *DiscordClient) GetAllRelationships(ctx context.Context) ([]Relationship, error) {
	body, err := dc.Request(ctx, "GET", "/users/@me/relationships")
	if err != nil {
		return nil, fmt.Errorf("error fetching relationships: %w", err)
	}
	defer body.Close()

	relationships := make([]Relationship, 0)

	if err := json.NewDecoder(body).Decode(&relationships); err != nil {
		return nil, fmt.Errorf("error parsing JSON response: %w", err)
	}

	return relationships, nil
}

// GetUserChannels retrieves all channels for the authenticated user
func (dc *DiscordClient) GetUserChannels(ctx context.Context) ([]Channel, error) {
	body, err := dc.Request(ctx, "GET", "/users/@me/channels")
	if err != nil {
		return nil, fmt.Errorf("error fetching channels: %w", err)
	}
	defer body.Close()

	channels := make([]Channel, 0)

	if err := json.NewDecoder(body).Decode(&channels); err != nil {
		return nil, fmt.Errorf("error parsing JSON response: %w", err)
	}

	return channels, nil
}

// CreateDMChannel creates a DM channel with the specified user
func (dc *DiscordClient) CreateDMChannel(ctx context.Context, userID string) (Channel, error) {
	payload := io.NopCloser(strings.NewReader(fmt.Sprintf(`{"recipient_id": "%s"}`, userID)))

	body, err := dc.RequestWithOptions(ctx, "POST", "/users/@me/channels", nil, payload)
	if err != nil {
		return Channel{}, fmt.Errorf("error creating DM channel: %w", err)
	}
	defer body.Close()

	var channel Channel
	if err := json.NewDecoder(body).Decode(&channel); err != nil {
		return Channel{}, fmt.Errorf("error parsing JSON response: %w", err)
	}

	return channel, nil
}

// RemoveDMChannel deletes a DM channel by its ID
func (dc *DiscordClient) RemoveDMChannel(ctx context.Context, channelID string) error {
	path := fmt.Sprintf("/channels/%s", channelID)
	_, err := dc.Request(ctx, "DELETE", path)
	if err != nil {
		return fmt.Errorf("error deleting DM channel: %w", err)
	}
	return nil
}

// GetMessages retrieves messages from a channel, paginated by the 'before' parameter
func (dc *DiscordClient) GetMessages(ctx context.Context, channelID string, before string) ([]map[string]any, error) {
	path := fmt.Sprintf("/channels/%s/messages", channelID)
	queries := url.Values{
		"limit": []string{MessageLimit},
	}

	if before != "" {
		queries.Add("before", before)
	}

	body, err := dc.RequestWithOptions(ctx, "GET", path, queries, nil)
	if err != nil {
		return nil, fmt.Errorf("error fetching messages: %w", err)
	}
	defer body.Close()

	var messages []map[string]any

	if err := json.NewDecoder(body).Decode(&messages); err != nil {
		return nil, fmt.Errorf("error parsing JSON response: %w", err)
	}
	return messages, nil
}

// GetUserGuilds retrieves the list of guilds (servers) the authenticated user is in
func (dc *DiscordClient) GetUserGuilds(ctx context.Context) ([]Guild, error) {
	body, err := dc.Request(ctx, "GET", "/users/@me/guilds")
	if err != nil {
		return nil, fmt.Errorf("error fetching guilds: %w", err)
	}
	defer body.Close()

	guilds := make([]Guild, 0)
	if err := json.NewDecoder(body).Decode(&guilds); err != nil {
		return nil, fmt.Errorf("error parsing JSON response: %w", err)
	}
	return guilds, nil
}

// GetGuildChannels retrieves all channels for a specific guild (server)
func (dc *DiscordClient) GetGuildChannels(ctx context.Context, guildID string) ([]Channel, error) {
	path := fmt.Sprintf("/guilds/%s/channels", guildID)
	body, err := dc.Request(ctx, "GET", path)
	if err != nil {
		return nil, fmt.Errorf("error fetching guild channels: %w", err)
	}
	defer body.Close()

	channels := make([]Channel, 0)
	if err := json.NewDecoder(body).Decode(&channels); err != nil {
		return nil, fmt.Errorf("error parsing JSON response: %w", err)
	}
	return channels, nil
}
