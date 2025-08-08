package discord

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/term"
)

const (
	ApiVersion   = "v9"
	MessageLimit = "100" // Max is 100
)

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

	// Set standard headers
	req.Header = http.Header{
		"Authorization": []string{dc.token},
		"Accept":        []string{"application/json"},
	}

	if body != nil {
		req.Body = body
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
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

	return resp.Body, nil
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
		"limit": []string{MessageLimit}, // Adjust limit as needed
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
