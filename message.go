package disgopher

import (
	"encoding/json"
)

//ChannelMention ...
type ChannelMention struct {
	ChannelID   string `json:"id"`
	GuildID     string `json:"guild_id"`
	ChannelType int    `json:"type"`
	ChannelName string `json:"string"`
}

//Attachment ...
type Attachment struct {
	ID       string `json:"id"`
	Filename string `json:"filename"`
	Size     int    `json:"size"`
	URL      string `json:"url"`
	ProzyURL string `json:"proxy_url"`
	Height   string `json:"height"`
	Width    string `json:"width"`
}

//Message ...
type Message struct {
	state           *clientState
	Guild           *Guild       `json:"-"`
	Channel         interface{}  `json:"-"`
	Attachments     []Attachment `json:"attachments"`
	Author          *User        `json:"-"`
	Content         string       `json:"content"`
	EditedTimestamp string       `json:"edited_timestamp"`
	Embeds          []*Embed     `json:"embeds"`
	Flags           int          `json:"flags"`
	ID              string       `json:"id"`
	Member          *GuildMember `json:"-"`
	MentionEveryone bool         `json:"mention_everyone"`
	//MentionRoles []Role `json:"mention_roles"`
	MentionedChannelsRaw       []*ChannelMention   `json:"mention_channels"`
	MentionedGuildTextChannels []*GuildTextChannel `json:"-"`
	//Mentions []GuildMember? `json:"mentions"`
	Nonce  interface{} `json:"nonce"`
	Pinned bool        `json:"pinned"`
	//ReferencedMessage MessageReference
	Timestamp string `json:"timestamp"`
	TTS       bool   `json:"tts"`
	Type      int    `josn:"type"`
	//Reactions []Reaction
	WebhookID string `json:"webhook_id"`
	//Activity Activity `json:"activity"`
	//Application Application `json:"application"`
	//Stickers []Sticker `json:"stickers"`
	//ReferencedMessage Message
	ChannelID string `json:"channel_id"`
	GuildID   string `json:"guild_id"`
}

type messagePrivate struct {
	Author *User        `json:"author"`
	Member *GuildMember `json:"member"`
}

func mentionedGuildTextChannelsFactory(message *Message) {
	for index := range message.MentionedChannelsRaw {
		mention := message.MentionedChannelsRaw[index]
		if mention.ChannelType == ChannelType.GuildText {
			channel := message.state.GuildTextChannels[mention.ChannelID]
			if channel != nil {
				message.MentionedGuildTextChannels = append(message.MentionedGuildTextChannels, channel)
			}
		}
	}
}

func newMessage(state *clientState, channelID string, guildID string, data []byte) *Message {
	message := &Message{state: state}
	json.Unmarshal(data, message)
	if channelID != "" {
		message.Channel = state.GuildTextChannels[channelID]
	} else {
		message.Channel = state.GuildTextChannels[message.ChannelID]
	}
	if guildID != "" {
		message.Guild = state.Guilds[guildID]
	} else {
		message.Guild = state.Guilds[message.GuildID]
	}
	private := new(messagePrivate)
	json.Unmarshal(data, private)
	if message.WebhookID == "" {
		message.Author = state.Users[private.Author.ID]
		message.Member = message.Guild.Members[private.Author.ID]
		if message.Author == nil {
			userData, _ := json.Marshal(private.Author)
			message.Author = newUser(state, userData)
		}
		if message.Member == nil && message.Guild != nil {
			memberData, _ := json.Marshal(private.Member)
			message.Member = newGuildMember(message.Author, message.Guild, memberData)
		}
	}
	mentionedGuildTextChannelsFactory(message)
	state.Messages[message.ID] = message
	return message
}

//EmbedField ...
type EmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

//EmbedAuthor ...
type EmbedAuthor struct {
	Name         string `json:"name,omitempty"`
	URL          string `json:"string,omitempty"`
	IconURL      string `json:"icon_url,omitempty"`
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}

//EmbedFooter ...
type EmbedFooter struct {
	Text         string `json:"text"`
	IconURL      string `json:"icon_url,omitempty"`
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}

//EmbedAttachment ...
type EmbedAttachment struct {
	URL      string `json:"url,omitempty"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Height   int    `json:"height,omitempty"`
	Witch    int    `json:"width,omitempty"`
}

//EmbedProvider ...
type EmbedProvider struct {
	Name string `json:"string,omitempty"`
	URL  string `json:"url,omitempty"`
}

//Embed ...
type Embed struct {
	Title       string           `json:"title,omitempty"`
	Type        string           `json:"type,omitempty"`
	Description string           `json:"description,omitempty"`
	URL         string           `json:"url,omitempty"`
	Timestamp   string           `json:"timestamp,omitempty"`
	Color       int              `json:"color,omitempty"`
	Footer      *EmbedFooter     `json:"footer,omitempty"`
	Image       *EmbedAttachment `json:"image,omitempty"`
	Thumbnail   *EmbedAttachment `json:"thumbnail,omitempty"`
	Video       *EmbedAttachment `json:"video,omitempty"`
	Provider    *EmbedProvider   `json:"provider,omitempty"`
	Author      *EmbedAuthor     `json:"author,omitempty"`
	Fields      []*EmbedField    `json:"fields,omitempty"`
}

//AddField ...
func (embed *Embed) AddField(name string, value string, inline bool) *EmbedField {
	field := &EmbedField{Name: name, Value: value, Inline: inline}
	embed.Fields = append(embed.Fields, field)
	return field
}
