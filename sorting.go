package main

import (
	"slices"
	"strings"
)

// getChannelSortName returns the name to use for sorting a channel
func getChannelSortName(channel Channel) string {
	switch channel.Type {
	case ChannelGroupDM:
		if channel.Name != "" {
			return channel.Name
		}
		return "Unnamed Group"
	case ChannelDM:
		// For private DMs, use the first recipient's name
		if len(channel.Recipients) > 0 {
			return getName(channel.Recipients[0])
		}
		return "Unnamed DM"
	default:
		return "Unknown"
	}
}

func sortChannels(channels []Channel) {
	// Sort channels alphabetically based on their sort name
	slices.SortFunc(channels, func(a, b Channel) int {
		nameA := getChannelSortName(a)
		nameB := getChannelSortName(b)
		return strings.Compare(strings.ToLower(nameA), strings.ToLower(nameB))
	})
}

func sortRelationships(relationships []Relationship) {
	// Sort relationships by global name (username)
	slices.SortFunc(relationships, func(a, b Relationship) int {
		nameA := getName(a.User)
		nameB := getName(b.User)
		return strings.Compare(strings.ToLower(nameA), strings.ToLower(nameB))
	})
}
