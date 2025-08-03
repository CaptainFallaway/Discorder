package main

import (
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
	MessageLimit = "100" // Max is 100.
)

// DiscordClient represents a client for interacting with the Discord API
type DiscordClient struct {
	token  string
	client *http.Client
}

func NewDiscordClient(token string) *DiscordClient {
	return &DiscordClient{
		token:  token,
		client: &http.Client{},
	}
}

// getRequest creates a new HTTP request with the necessary headers
// and returns it. This is used internally to avoid code duplication.
func (dc *DiscordClient) getRequest(method, path string) http.Request {
	return http.Request{
		Method: method,
		URL:    &url.URL{Scheme: "https", Host: "discord.com", Path: fmt.Sprintf("/api/%s%s", ApiVersion, path)},
		Header: http.Header{
			"Authorization": []string{dc.token},
		},
	}
}

// Simple request with just method and path
func (dc *DiscordClient) Request(method, path string) (io.ReadCloser, error) {
	return dc.RequestWithOptions(method, path, nil, nil)
}

// Full request with queries and body
func (dc *DiscordClient) RequestWithOptions(method, path string, queries url.Values, body io.ReadCloser) (io.ReadCloser, error) {
	request := dc.getRequest(method, path)

	if body != nil {
		request.Body = body
		request.Header.Set("Content-Type", "application/json")
	}

	if queries != nil {
		request.URL.RawQuery = queries.Encode()
	}

	// Only print debug info if output is going to a terminal (not piped)
	if term.IsTerminal(int(os.Stdout.Fd())) {
		fmt.Println("Making request:", request.Method, request.URL.String())
	}

	resp, err := dc.client.Do(&request)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("request failed with status: %s", resp.Status)
	}

	return resp.Body, nil
}

// GetAllRelationships retrieves all relationships for the authenticated user
func (dc *DiscordClient) GetAllRelationships() ([]Relationship, error) {
	body, err := dc.Request("GET", "/users/@me/relationships")
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
func (dc *DiscordClient) GetUserChannels() ([]Channel, error) {
	body, err := dc.Request("GET", "/users/@me/channels")
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
func (dc *DiscordClient) CreateDMChannel(userID string) (Channel, error) {
	payload := io.NopCloser(strings.NewReader(fmt.Sprintf(`{"recipient_id": "%s"}`, userID)))

	body, err := dc.RequestWithOptions("POST", "/users/@me/channels", nil, payload)
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
func (dc *DiscordClient) RemoveDMChannel(channelID string) error {
	path := fmt.Sprintf("/channels/%s", channelID)
	_, err := dc.Request("DELETE", path)
	if err != nil {
		return fmt.Errorf("error deleting DM channel: %w", err)
	}
	return nil
}

// GetMessages retrieves messages from a channel, paginated by the 'before' parameter
func (dc *DiscordClient) GetMessages(channelID string, before string) ([]map[string]any, error) {
	path := fmt.Sprintf("/channels/%s/messages", channelID)
	queries := url.Values{
		"limit": []string{MessageLimit}, // Adjust limit as needed
	}

	if before != "" {
		queries.Add("before", before)
	}

	body, err := dc.RequestWithOptions("GET", path, queries, nil)
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
