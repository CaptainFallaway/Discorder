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
		return fmt.Errorf("no relationships found")
	}

	sortRelationships(relationships)

	// Create table data
	tableData := [][]string{
		{"User ID", "Global Name (Username) aka [Nickname]", "Type", "Since"},
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
		since = fmt.Sprintf("%s (%s)", formatTime(rel.Since), since)

		if rel.Nickname != "" {
			name = fmt.Sprintf("%s aka [%s]", name, rel.Nickname)
		}

		tableData = append(tableData, []string{
			rel.User.ID,
			name,
			relationshipType,
			since,
		})
	}

	fmt.Printf("Found %d relationships:\n\n", len(relationships))
	pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()

	return nil
}
