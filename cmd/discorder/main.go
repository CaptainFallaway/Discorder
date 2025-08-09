package main

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/term"

	"github.com/CaptainFallaway/Discorder/internal/cli"
	"github.com/CaptainFallaway/Discorder/internal/discord"
)

func run(token, action string, args []string) error {
	// Only print debug info if output is going to a terminal (not piped)
	debug := term.IsTerminal(int(os.Stdout.Fd()))

	dc := discord.NewDiscordClient(token, debug)

	switch action {
	case "relationships":
		if err := cli.PrintRelationships(dc); err != nil {
			return fmt.Errorf("error printing relationships: %w", err)
		}
	case "dms":
		if err := cli.PrintDMs(dc); err != nil {
			return fmt.Errorf("error printing DMs: %w", err)
		}
	case "create-dm":
		if len(args) < 1 {
			return fmt.Errorf("user ID is required to create / retrieve a DM channel")
		}
		userID := args[0]
		channel, err := dc.CreateDMChannel(context.Background(), userID)
		if err != nil {
			return fmt.Errorf("error creating DM channel: %w", err)
		}
		fmt.Printf("DM channel created with ID: %s\n", channel.ID)
	case "remove-dm":
		if len(args) < 1 {
			return fmt.Errorf("channel ID is required to delete a DM channel")
		}
		channelID := args[0]
		if err := dc.RemoveDMChannel(context.Background(), channelID); err != nil {
			return fmt.Errorf("error deleting DM channel: %w", err)
		}
		fmt.Printf("DM channel with ID %s deleted successfully.\n", channelID)
	case "guilds":
		if err := cli.PrintGuilds(dc); err != nil {
			return fmt.Errorf("error printing guilds: %w", err)
		}
	case "guild-channels":
		if len(args) < 1 {
			return fmt.Errorf("guild ID is required to list channels")
		}
		guildID := args[0]
		if err := cli.PrintGuildChannels(dc, guildID); err != nil {
			return fmt.Errorf("error printing guild channels: %w", err)
		}
	case "messages":
		if len(args) < 1 {
			return fmt.Errorf("channel ID is required to dump messages")
		}
		channelID := args[0]
		messages, err := cli.GetAllMessages(dc, channelID)
		if err != nil {
			return fmt.Errorf("error fetching messages: %w", err)
		}
		cli.PrettyPrintJSON(messages)
	default:
		fmt.Printf("Unknown action \"%s\". Available actions: relationships, dms, create-dm, remove-dm, guilds, guild-channels, messages\n", action)
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
			fmt.Println("Usage: ./discorder <action> [args...]")
			fmt.Println("Available actions: relationships, dms, create-dm, remove-dm, guilds, guild-channels, messages")
			os.Exit(1)
		}

		action = os.Args[1]
		args = os.Args[2:]
	} else {
		// Without env token, we need at least 3 args: program_name, token, and action
		if len(os.Args) < 3 {
			fmt.Println("Must provide a Discord Token and an action")
			fmt.Println("Usage: ./discorder <token> <action> [args...]")
			fmt.Println("   or: DISCORD_TOKEN=your_token ./discorder <action> [args...]")
			fmt.Println("Available actions: relationships, dms, create-dm, remove-dm, guilds, guild-channels, messages")
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
