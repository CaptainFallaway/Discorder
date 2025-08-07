package main

import (
	"fmt"
	"slices"
)

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

		if len(messages) < 100 {
			break
		}
	}

	slices.Reverse(allMessages)

	return allMessages, nil
}
