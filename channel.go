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

//ChannelBase ...
type ChannelBase struct {
	state *clientState
	data  []byte
	ID    string `json:"id"`
	Type  int    `json:"type"`
}

func (channel *ChannelBase) upgrade(args ...interface{}) interface{} {
	switch channel.Type {
	case ChannelType.GuildText:
		return newGuildTextChannel(channel, args[0].(*Guild), channel.data)
	default:
		return nil
	}
}

//GuildTextChannel ...
type GuildTextChannel struct {
	state    *clientState
	Guild    *Guild
	Type     int
	Topic    string `json:"topic"`
	Cooldown int    `json:"ratelimit_per_user"`
	Position int    `json:"position"`
	//PremissionOverwrites []PermissionOverwrite `json:"permission_overwrites"`
	ParentID      string `json:"parent_id"`
	NSFW          bool   `json:"nsfw"`
	Name          string `json:"name"`
	LastMessageID string `json:"last_message_id"`
	ID            string
	GuildID       string `json:"guild_id"`
}

func newGuildTextChannel(baseChannel *ChannelBase, guild *Guild, data []byte) *GuildTextChannel {
	channel := &GuildTextChannel{state: baseChannel.state, ID: baseChannel.ID, Type: baseChannel.Type}
	json.Unmarshal(data, channel)
	if guild == nil {
		if channel.GuildID != "" {
			gotGuild := channel.state.Guilds[channel.GuildID]
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
