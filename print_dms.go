package main

import (
	"fmt"

	"github.com/pterm/pterm"
)

func printDms(dc *DiscordClient) error {
	channels, err := dc.GetUserChannels()
	if err != nil {
		return fmt.Errorf("failed to get user channels: %w", err)
	}

	// Filter only DM channels
	dmChannels := make([]Channel, 0)
	for _, channel := range channels {
		if channel.Type == ChannelDM || channel.Type == ChannelGroupDM {
			dmChannels = append(dmChannels, channel)
		}
	}

	if len(dmChannels) == 0 {
		fmt.Println("No direct message channels found.")
		return fmt.Errorf("no direct message channels found")
	}

	// Create table data
	tableData := [][]string{
		{"Channel ID", "Type", "Name", "Recipients"},
	}

	for _, channel := range dmChannels {
		channelType := ""
		switch channel.Type {
		case ChannelDM:
			channelType = "DM"
		case ChannelGroupDM:
			channelType = "Group DM"
		default:
			channelType = "Unknown"
		}

		// Build recipients list
		recipients := ""
		for i, recipient := range channel.Recipients {
			// name := getName(recipient)
			name := recipient.Username

			if i > 0 {
				recipients += ", "
			}
			recipients += name
		}

		// Use channel name if available (for group DMs), otherwise use recipients
		channelName := channel.Name
		if channelName == "" {
			channelName = recipients
		}

		tableData = append(tableData, []string{
			channel.ID,
			channelType,
			channelName,
			recipients,
		})
	}

	fmt.Printf("Found %d direct message channels:\n\n", len(dmChannels))
	pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()

	return nil
}
