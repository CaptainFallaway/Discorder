package cli

import (
	"slices"
	"strings"

	"github.com/CaptainFallaway/Discorder/internal/discord"
)

// getChannelSortName returns the name to use for sorting a channel
func getChannelSortName(channel discord.Channel) string {
	switch channel.Type {
	case discord.ChannelGroupDM:
		if channel.Name != "" {
			return channel.Name
		}
		return "Unnamed Group"
	case discord.ChannelDM:
		// For private DMs, use the first recipient's name
		if len(channel.Recipients) > 0 {
			return channel.Recipients[0].GetName()
		}
		return "Unnamed DM"
	default:
		return "Unknown"
	}
}

func SortChannels(channels []discord.Channel) {
	// Sort channels alphabetically based on their sort name
	slices.SortFunc(channels, func(a, b discord.Channel) int {
		nameA := getChannelSortName(a)
		nameB := getChannelSortName(b)
		return strings.Compare(strings.ToLower(nameA), strings.ToLower(nameB))
	})
}

func SortRelationships(relationships []discord.Relationship) {
	// Sort relationships by global name (username)
	slices.SortFunc(relationships, func(a, b discord.Relationship) int {
		nameA := a.User.GetName()
		nameB := b.User.GetName()
		return strings.Compare(strings.ToLower(nameA), strings.ToLower(nameB))
	})
}
