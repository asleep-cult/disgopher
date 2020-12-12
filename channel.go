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
		if len(args) > 0 {
			if guild, ok := args[0].(*Guild); ok {
				return newGuildTextChannel(channel, guild, channel.data)
			}
		}
		return newGuildTextChannel(channel, nil, channel.data)
	case ChannelType.GuildVoice:
		if len(args) > 0 {
			if guild, ok := args[0].(*Guild); ok {
				return newGuildVoiceChannel(channel, guild, channel.data)
			}
		}
		return newGuildVoiceChannel(channel, nil, channel.data)
	default:
		return nil
	}
}

func newBaseChannel(state *clientState, data []byte) *ChannelBase {
	channel := &ChannelBase{state: state, data: data}
	json.Unmarshal(data, channel)
	return channel
}

//GuildTextChannel ...
type GuildTextChannel struct {
	state    *clientState
	Guild    *Guild `json:"-"`
	Type     int    `json:"type"`
	Topic    string `json:"topic"`
	SlowMode int    `json:"ratelimit_per_user"`
	Position int    `json:"position"`
	//PremissionOverwrites []PermissionOverwrite `json:"permission_overwrites"`
	ParentID      string   `json:"parent_id"`
	NSFW          bool     `json:"nsfw"`
	Name          string   `json:"name"`
	LastMessage   *Message `json:"-"`
	LastMessageID string   `json:"last_message_id"`
	ID            string   `json:"id"`
	GuildID       string   `json:"guild_id"`
}

//Send ...
func (channel *GuildTextChannel) Send(req *MessageCreateRequest) (*Message, error) {
	data, err := channel.state.HTTP.messageCreate(channel.ID, req)
	var message *Message
	if err == nil {
		message = newMessage(channel.state, channel.ID, channel.Guild.ID, data)
	}
	return message, err
}

type channelNameUpdateReq struct {
	Name string `json:"name"`
}

//SetName ...
func (channel *GuildTextChannel) SetName(name string) error {
	data, err := channel.state.HTTP.modifyChannel(
		channel.ID,
		channelNameUpdateReq{name})
	if err == nil {
		channel.update(data)
	}
	return err
}

type channelTypeUpdateReq struct {
	Type int `json:"type"`
}

//SetType ...
func (channel *GuildTextChannel) SetType(channelType int) error {
	data, err := channel.state.HTTP.modifyChannel(
		channel.ID,
		channelTypeUpdateReq{channelType})
	if err == nil {
		channel.update(data)
	}
	return err
}

type channelPositionUpdateReq struct {
	Position int `json:"position"`
}

//SetPosition ...
func (channel *GuildTextChannel) SetPosition(position int) error {
	data, err := channel.state.HTTP.modifyChannel(
		channel.ID,
		channelPositionUpdateReq{position})
	if err == nil {
		channel.update(data)
	}
	return err
}

type channelTopicUpdateReq struct {
	Topic string `json:"topic"`
}

//SetTopic ...
func (channel *GuildTextChannel) SetTopic(topic string) error {
	data, err := channel.state.HTTP.modifyChannel(
		channel.ID,
		channelTopicUpdateReq{topic})
	if err == nil {
		channel.update(data)
	}

	return err
}

type channelNSFWUpdateReq struct {
	NSFW bool `json:"nsfw"`
}

//SetNSFW ...
func (channel *GuildTextChannel) SetNSFW(nsfw bool) error {
	data, err := channel.state.HTTP.modifyChannel(
		channel.ID,
		channelNSFWUpdateReq{nsfw})
	if err == nil {
		channel.update(data)
	}
	return err
}

type channelSlowModeUpdateReq struct {
	SlowMode int `json:"ratelimit_per_user"`
}

//SetSlowMode ...
func (channel *GuildTextChannel) SetSlowMode(cooldown int) error {
	data, err := channel.state.HTTP.modifyChannel(
		channel.ID,
		channelSlowModeUpdateReq{cooldown})
	if err == nil {
		channel.update(data)
	}
	return err
}

func (channel *GuildTextChannel) update(data []byte) {
	json.Unmarshal(data, channel)
}

func newGuildTextChannel(baseChannel *ChannelBase, guild *Guild, data []byte) *GuildTextChannel {
	channel := &GuildTextChannel{state: baseChannel.state}
	json.Unmarshal(data, channel)
	if guild != nil {
		channel.Guild = guild
	} else {
		channel.Guild = channel.state.Guilds[channel.GuildID]
	}
	channel.LastMessage = channel.state.Messages[channel.LastMessageID]
	channel.state.GuildTextChannels[channel.ID] = channel
	channel.Guild.TextChannels[channel.ID] = channel
	return channel
}

//GuildVoiceChannel ...
type GuildVoiceChannel struct {
	state    *clientState
	Guild    *Guild `json:"-"`
	Bitrate  int    `json:"bitrate"`
	GuildID  string `json:"guild_id"`
	ID       string `json:"id"`
	Name     string `json:"name"`
	NSFW     bool   `json:"nsfw"`
	ParentID string `json:"parent_id"`
	//PermissionOverwrites []PermissionOverwrite
	Position  int `json:"position"`
	Type      int
	UserLimit int `json:"user_limit"`
}

func newGuildVoiceChannel(baseChannel *ChannelBase, guild *Guild, data []byte) *GuildVoiceChannel {
	channel := &GuildVoiceChannel{state: baseChannel.state, ID: baseChannel.ID, Type: baseChannel.Type}
	json.Unmarshal(data, channel)
	if guild != nil {
		channel.Guild = guild
	} else {
		channel.Guild = channel.state.Guilds[channel.GuildID]
	}
	channel.state.GuildVoiceChannels[channel.ID] = channel
	channel.Guild.VoiceChannels[channel.ID] = channel
	return channel
}

func (channel *GuildVoiceChannel) update(data []byte) {
	json.Unmarshal(data, channel)
}
