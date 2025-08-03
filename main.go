package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func run(token, action string, args []string) error {
	dc := NewDiscordClient(token)

	switch action {
	case "rels":
		if err := printRelationships(dc); err != nil {
			return fmt.Errorf("error printing relationships: %w", err)
		}
	case "dms":
		if err := printDms(dc); err != nil {
			return fmt.Errorf("error printing DMs: %w", err)
		}
	case "gdm":
		if len(args) < 1 {
			return fmt.Errorf("user ID is required to create / retrieve a DM channel")
		}
		userID := args[0]
		channel, err := dc.GetDMChannel(userID)
		if err != nil {
			return fmt.Errorf("error creating DM channel: %w", err)
		}
		fmt.Printf("DM channel created with ID: %s\n", channel.ID)
	case "rdm":
		if len(args) < 1 {
			return fmt.Errorf("channel ID is required to delete a DM channel")
		}
		channelID := args[0]
		if err := dc.RemoveDMChannel(channelID); err != nil {
			return fmt.Errorf("error deleting DM channel: %w", err)
		}
		fmt.Printf("DM channel with ID %s deleted successfully.\n", channelID)
	case "msgs":
		if len(args) < 1 {
			return fmt.Errorf("channel ID is required to dump messages")
		}
		channelID := args[0]
		messages, err := getAllMessages(dc, channelID)
		if err != nil {
			return fmt.Errorf("error fetching messages: %w", err)
		}
		prettyPrintJson(messages)
	default:
		fmt.Printf("Unknown action \"%s\". Available actions: rels, dms, gdm, rdm, msgs\n", action)
	}

	return nil
}

func main() {
	godotenv.Load()

	var token string
	var action string
	var args []string

	// Check if token is provided via environment variable
	if envToken := os.Getenv("DISCORD_TOKEN"); envToken != "" {
		token = envToken

		// With env token, we need at least 2 args: program_name and action
		if len(os.Args) < 2 {
			fmt.Println("Must provide an action when using the DISCORD_TOKEN environment variable")
			fmt.Println("Usage: ./program <action> [args...]")
			fmt.Println("Available actions: rels, dms, gdm, rdm, msgs")
			os.Exit(1)
		}

		action = os.Args[1]
		args = os.Args[2:]
	} else {
		// Without env token, we need at least 3 args: program_name, token, and action
		if len(os.Args) < 3 {
			fmt.Println("Must provide a Discord Token and an action")
			fmt.Println("Usage: ./program <token> <action> [args...]")
			fmt.Println("   or: DISCORD_TOKEN=your_token ./program <action> [args...]")
			fmt.Println("Available actions: rels, dms, gdm, rdm, msgs")
			os.Exit(1)
		}

		token = os.Args[1]
		action = os.Args[2]
		args = os.Args[3:]
	}

	if err := run(token, action, args); err != nil {
		fmt.Println(err)
	}
}
