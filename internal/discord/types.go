package discord

import "fmt"

// Relationship types
const (
	RelationNone = iota
	RelationFriend
	RelationBlocked
	RelationPendingIncoming
	RelationPendingOutgoing
	RelationImplicit
)

// Channel types
const (
	ChannelText               = 0
	ChannelDM                 = 1
	ChannelVoice              = 2
	ChannelGroupDM            = 3
	ChannelGuildCategory      = 4
	ChannelGuildAnnouncement  = 5
	ChannelAnnouncementThread = 10
	ChannelPublicThread       = 11
	ChannelPrivateThread      = 12
	ChannelGuildStageVoice    = 13
	ChannelGuildDirectory     = 14
	ChannelGuildForum         = 15
	ChannelGuildMedia         = 16
)

// User is a partial Discord user.
type User struct {
	ID         string `json:"id"`
	Username   string `json:"username"`
	GlobalName string `json:"global_name"`
	Avatar     string `json:"avatar"`
}

// GetName returns the user's name, preferring global name if available.
func (u User) GetName() string {
	if u.GlobalName == "" {
		return u.Username
	}
	return fmt.Sprintf("%s (%s)", u.GlobalName, u.Username)
}

// Relationship is a partial relationship record.
type Relationship struct {
	ID       string `json:"id"`
	Type     int    `json:"type"`
	Nickname string `json:"nickname"`
	User     User   `json:"user"`
	Since    string `json:"since"`
}

// Channel is a partial channel object.
type Channel struct {
	ID         string `json:"id"`
	Type       int    `json:"type"`
	Name       string `json:"name"`
	Recipients []User `json:"recipients"`
	NSFW       bool   `json:"nsfw"`
}

// Guild NSFW levels
const (
	NSFWLevelDefault       = 0
	NSFWLevelExplicit      = 1
	NSFWLevelSafe          = 2
	NSFWLevelAgeRestricted = 3
)

// Guild represents a partial guild object as returned by /users/@me/guilds.
type Guild struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Owner       bool   `json:"owner"`
	NSFWLevel   int    `json:"nsfw_level"`
	Description string `json:"description"`
}
