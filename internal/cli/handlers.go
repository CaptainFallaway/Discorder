package cli

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/pterm/pterm"

	"github.com/CaptainFallaway/Discorder/internal/discord"
)

func PrintDMs(dc *discord.DiscordClient) error {
	channels, err := dc.GetUserChannels(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get user channels: %w", err)
	}

	SortChannels(channels)

	// Separate channels by type
	groupDMs := make([]discord.Channel, 0)
	privateDMs := make([]discord.Channel, 0)

	for _, channel := range channels {
		switch channel.Type {
		case discord.ChannelGroupDM:
			groupDMs = append(groupDMs, channel)
		case discord.ChannelDM:
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

func printGroupDMs(channels []discord.Channel) error {
	if len(channels) == 0 {
		fmt.Println("No group DM channels found.")
		return nil
	}

	// Create table data for Group DMs
	tableData := [][]string{{"Channel ID", "Name", "Recipients"}}

	for _, channel := range channels {
		// Build recipients list
		recipients := ""
		for i, recipient := range channel.Recipients {
			name := recipient.GetName()
			if i > 0 {
				recipients += ", "
			}
			recipients += name
		}

		channelName := channel.Name
		if channelName == "" {
			channelName = "Unnamed Group"
		}

		tableData = append(tableData, []string{channel.ID, channelName, recipients})
	}

	fmt.Printf("Found %d group DM channels:\n\n", len(channels))
	pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
	fmt.Println()

	return nil
}

func printPrivateDMs(channels []discord.Channel) error {
	if len(channels) == 0 {
		fmt.Println("No private DM channels found.")
		return nil
	}

	// Create table data for Private DMs (no recipients column)
	tableData := [][]string{{"Channel ID", "User"}}

	for _, channel := range channels {
		// For private DMs, show the other user's name
		userName := "Unknown User"
		if len(channel.Recipients) > 0 {
			userName = channel.Recipients[0].GetName()
		}

		tableData = append(tableData, []string{channel.ID, userName})
	}

	fmt.Printf("Found %d private DM channels:\n\n", len(channels))
	pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()

	return nil
}

func PrintRelationships(dc *discord.DiscordClient) error {
	relationships, err := dc.GetAllRelationships(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get relationships: %w", err)
	}

	if len(relationships) == 0 {
		return fmt.Errorf("no relationships found")
	}

	SortRelationships(relationships)

	// Create table data
	tableData := [][]string{{"User ID", "Global Name (Username) aka [Nickname]", "Type", "Since"}}

	for _, rel := range relationships {
		relationshipType := relationshipTypeString(rel.Type)

		name := rel.User.GetName()
		since := FormatTimeSince(rel.Since)
		since = fmt.Sprintf("%s (%s)", FormatTime(rel.Since), since)

		if rel.Nickname != "" {
			name = fmt.Sprintf("%s aka [%s]", name, rel.Nickname)
		}

		tableData = append(tableData, []string{rel.User.ID, name, relationshipType, since})
	}

	fmt.Printf("Found %d relationships:\n\n", len(relationships))
	pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()

	return nil
}

func relationshipTypeString(t int) string {
	switch t {
	case discord.RelationFriend:
		return "Friend"
	case discord.RelationBlocked:
		return "Blocked"
	case discord.RelationPendingIncoming:
		return "Pending Incoming"
	case discord.RelationPendingOutgoing:
		return "Pending Outgoing"
	case discord.RelationImplicit:
		return "Implicit"
	default:
		return "Unknown"
	}
}

func PrintGuilds(dc *discord.DiscordClient) error {
	guilds, err := dc.GetUserGuilds(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get guilds: %w", err)
	}

	if len(guilds) == 0 {
		fmt.Println("No guilds found.")
		return nil
	}

	// sort by name
	slices.SortFunc(guilds, func(a, b discord.Guild) int {
		return strings.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name))
	})

	table := [][]string{{"Guild ID", "Name", "Owner Of", "NSFW Level", "Description"}}
	for _, g := range guilds {
		owner := "No"
		if g.Owner {
			owner = "Yes"
		}
		nsfw := nsfwLevelString(g.NSFWLevel)
		desc := g.Description
		if desc == "" {
			desc = "-"
		}
		table = append(table, []string{g.ID, g.Name, owner, nsfw, desc})
	}

	fmt.Printf("Found %d guilds:\n\n", len(guilds))
	pterm.DefaultTable.WithHasHeader().WithData(table).Render()
	return nil
}

func nsfwLevelString(level int) string {
	switch level {
	case discord.NSFWLevelDefault:
		return "Default"
	case discord.NSFWLevelExplicit:
		return "Explicit"
	case discord.NSFWLevelSafe:
		return "Safe"
	case discord.NSFWLevelAgeRestricted:
		return "Age Restricted"
	default:
		return fmt.Sprintf("Unknown(%d)", level)
	}
}

func PrintGuildChannels(dc *discord.DiscordClient, guildID string) error {
	if guildID == "" {
		return fmt.Errorf("guild ID is required")
	}

	chns, err := dc.GetGuildChannels(context.Background(), guildID)
	if err != nil {
		return fmt.Errorf("failed to get guild channels: %w", err)
	}

	if len(chns) == 0 {
		fmt.Println("No channels found.")
		return nil
	}

	// sort channels: by type then by name
	slices.SortFunc(chns, func(a, b discord.Channel) int {
		if a.Type != b.Type {
			return a.Type - b.Type
		}
		return strings.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name))
	})

	table := [][]string{{"Channel ID", "Type", "Name", "NSFW"}}
	for _, c := range chns {
		nsfw := "No"
		if c.NSFW {
			nsfw = "Yes"
		}
		table = append(table, []string{c.ID, channelTypeString(c.Type), c.Name, nsfw})
	}

	fmt.Printf("Found %d channels:\n\n", len(chns))
	pterm.DefaultTable.WithHasHeader().WithData(table).Render()
	return nil
}

// channelTypeString returns a human-readable name for a channel type
func channelTypeString(t int) string {
	switch t {
	case discord.ChannelText:
		return "Text"
	case discord.ChannelDM:
		return "DM"
	case discord.ChannelVoice:
		return "Voice"
	case discord.ChannelGroupDM:
		return "Group DM"
	case discord.ChannelGuildCategory:
		return "Guild Category"
	case discord.ChannelGuildAnnouncement:
		return "Guild Announcement"
	case discord.ChannelAnnouncementThread:
		return "Announcement Thread"
	case discord.ChannelPublicThread:
		return "Public Thread"
	case discord.ChannelPrivateThread:
		return "Private Thread"
	case discord.ChannelGuildStageVoice:
		return "Stage Voice"
	case discord.ChannelGuildDirectory:
		return "Guild Directory"
	case discord.ChannelGuildForum:
		return "Guild Forum"
	case discord.ChannelGuildMedia:
		return "Guild Media"
	default:
		return fmt.Sprintf("Unknown(%d)", t)
	}
}
