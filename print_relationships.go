package main

import (
	"fmt"

	"github.com/pterm/pterm"
)

func printRelationships(dc *DiscordClient) error {
	relationships, err := dc.GetAllRelationships()
	if err != nil {
		return fmt.Errorf("failed to get relationships: %w", err)
	}

	if len(relationships) == 0 {
		fmt.Println("No relationships found.")
		return fmt.Errorf("no relationships found")
	}

	// Create table data
	tableData := [][]string{
		{"Global Name (Username)", "Nickname", "User ID", "Type", "Since"},
	}

	for _, rel := range relationships {
		relationshipType := ""
		switch rel.Type {
		case RelationFriend:
			relationshipType = "Friend"
		case RelationBlocked:
			relationshipType = "Blocked"
		case RelationPendingIncoming:
			relationshipType = "Pending Incoming"
		case RelationPendingOutgoing:
			relationshipType = "Pending Outgoing"
		case RelationImplicit:
			relationshipType = "Implicit"
		default:
			relationshipType = "Unknown"
		}

		name := getName(rel.User)
		since := formatTimeSince(rel.Since)

		tableData = append(tableData, []string{
			name,
			rel.Nickname,
			rel.User.ID,
			relationshipType,
			since,
		})
	}

	fmt.Printf("Found %d relationships:\n\n", len(relationships))
	pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()

	return nil
}
