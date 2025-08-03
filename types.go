package main

const (
	RelationNone = iota
	RelationFriend
	RelationBlocked
	RelationPendingIncoming
	RelationPendingOutgoing
	RelationImplicit
)

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

type User struct {
	ID         string `json:"id"`
	Username   string `json:"username"`
	GlobalName string `json:"global_name"`
	Avatar     string `json:"avatar"`
}

type Relationship struct {
	ID       string `json:"id"`
	Type     int    `json:"type"`
	Nickname string `json:"nickname"`
	User     User   `json:"user"`
	Since    string `json:"since"`
}

type Channel struct {
	ID         string `json:"id"`
	Type       int    `json:"type"`
	Name       string `json:"name"`
	Recipients []User `json:"recipients"`
}
