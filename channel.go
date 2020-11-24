package disgopher

import (
	"encoding/json"
)

//ChannelType ...
var ChannelType = struct {
	GuildText     int
	DM            int
	GuildVoice    int
	GroupDM       int
	GuildCaregory int
	GuildNews     int
	GuildStore    int
}{0, 1, 2, 3, 4, 5, 6}

//GuildTextChannel ...
type GuildTextChannel struct {
	Guild    *Guild
	Type     int    `json:"type"`
	Topic    string `json:"topic"`
	Cooldown int    `json:"ratelimit_per_user"`
	Position int    `json:"position"`
	//PremissionOverwrites []PermissionOverwrite `json:"permission_overwrites"`
	ParentID      string `json:"parent_id"`
	NSFW          bool   `json:"nsfw"`
	Name          string `json:"name"`
	LastMessageID string `json:"last_message_id"`
	ID            string `json:"id"`
	GuildID       string `json:"guild_id"`
}

func newGuildTextChannel(state *clientState, guild *Guild, data []byte) *GuildTextChannel {
	channel := new(GuildTextChannel)
	json.Unmarshal(data, channel)
	if guild == nil {
		if channel.GuildID != "" {
			gotGuild := state.Guilds[channel.GuildID]
			if gotGuild != nil {
				channel.Guild = gotGuild
			}
		}
	} else {
		channel.Guild = guild
	}
	channel.Guild = guild
	return channel
}
