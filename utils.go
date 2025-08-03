package main

import (
	"fmt"
	"slices"
	"time"

	"github.com/hokaccha/go-prettyjson"
)

func formatTimeSince(sinceStr string) string {
	if sinceStr == "" {
		return "Unknown"
	}

	// Parse the ISO 8601 date string
	since, err := time.Parse(time.RFC3339, sinceStr)
	if err != nil {
		return sinceStr // Return original if parsing fails
	}

	duration := time.Since(since)
	days := int(duration.Hours() / 24)

	if days < 1 {
		hours := int(duration.Hours())
		if hours < 1 {
			return "Less than an hour ago"
		}
		return fmt.Sprintf("%d hour(s) ago", hours)
	} else if days < 30 {
		return fmt.Sprintf("%d day(s) ago", days)
	} else if days < 365 {
		months := days / 30
		return fmt.Sprintf("%d month(s) ago", months)
	} else {
		years := days / 365
		remainingDays := days % 365
		if remainingDays < 30 {
			return fmt.Sprintf("%d year(s) ago", years)
		} else {
			months := remainingDays / 30
			return fmt.Sprintf("%d year(s), %d month(s) ago", years, months)
		}
	}

}

func getName(user User) string {
	if user.GlobalName == "" {
		return user.Username
	}
	return fmt.Sprintf("%s (%s)", user.GlobalName, user.Username)
}

func prettyPrintJson(v any) {
	prettyjson, err := prettyjson.Marshal(v)
	if err != nil {
		fmt.Println("Error marshaling to pretty JSON:", err)
		return
	}
	fmt.Println(string(prettyjson))
}

func getAllMessages(dc *DiscordClient, channelID string) ([]map[string]any, error) {
	allMessages := make([]map[string]any, 0, 100)
	var before string

	for {
		messages, err := dc.GetMessages(channelID, before)
		if err != nil {
			return nil, fmt.Errorf("error fetching messages: %w", err)
		}

		if len(messages) == 0 {
			break
		}

		allMessages = append(allMessages, messages...)
		before = messages[len(messages)-1]["id"].(string) // Use the last message's ID for the next request
	}

	slices.Reverse(allMessages)

	return allMessages, nil
}
