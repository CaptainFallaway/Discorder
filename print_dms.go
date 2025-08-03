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

	sortChannels(channels)

	// Separate channels by type
	groupDMs := make([]Channel, 0)
	privateDMs := make([]Channel, 0)

	for _, channel := range channels {
		switch channel.Type {
		case ChannelGroupDM:
			groupDMs = append(groupDMs, channel)
		case ChannelDM:
			privateDMs = append(privateDMs, channel)
		}
	}

	// Print Group DMs first
	if err := printGroupDMs(groupDMs); err != nil {
		return err
	}

	// Print Private DMs second
	if err := printPrivateDMs(privateDMs); err != nil {
		return err
	}

	return nil
}

func printGroupDMs(channels []Channel) error {
	if len(channels) == 0 {
		fmt.Println("No group DM channels found.")
		return nil
	}

	// Create table data for Group DMs
	tableData := [][]string{
		{"Channel ID", "Name", "Recipients"},
	}

	for _, channel := range channels {
		// Build recipients list
		recipients := ""
		for i, recipient := range channel.Recipients {
			name := getName(recipient)
			if i > 0 {
				recipients += ", "
			}
			recipients += name
		}

		channelName := channel.Name
		if channelName == "" {
			channelName = "Unnamed Group"
		}

		tableData = append(tableData, []string{
			channel.ID,
			channelName,
			recipients,
		})
	}

	fmt.Printf("Found %d group DM channels:\n\n", len(channels))
	pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
	fmt.Println()

	return nil
}

func printPrivateDMs(channels []Channel) error {
	if len(channels) == 0 {
		fmt.Println("No private DM channels found.")
		return nil
	}

	// Create table data for Private DMs (no recipients column)
	tableData := [][]string{
		{"Channel ID", "User"},
	}

	for _, channel := range channels {
		// For private DMs, show the other user's name
		userName := "Unknown User"
		if len(channel.Recipients) > 0 {
			userName = getName(channel.Recipients[0])
		}

		tableData = append(tableData, []string{
			channel.ID,
			userName,
		})
	}

	fmt.Printf("Found %d private DM channels:\n\n", len(channels))
	pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()

	return nil
}
