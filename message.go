package disgopher

import "encoding/json"

//ChannelMention ...
type ChannelMention struct {
	ChannelID   string `json:"id"`
	GuildID     string `json:"guild_id"`
	ChannelType int    `json:"type"`
	ChannelName string `json:"string"`
}

//Message ...
type Message struct {
	state           *clientState
	Guild           *Guild
	Channel         *GuildTextChannel
	Attachments     []Attachment `json:"attachments"`
	Author          *User        `json:"-"`
	Content         string       `json:"content"`
	EditedTimestamp string       `json:"edited_timestamp"`
	//Embeds []Embed
	Flags int    `json:"flags"`
	ID    string `json:"id"`
	//Member *GuildMember
	MentionEveryone bool `json:"mention_everyone"`
	//MentionRoles []Role `json:"mention_roles"`
	MentionedChannelsRaw       []*ChannelMention `json:"mention_channels"`
	MentionedGuildTextChannels []*GuildTextChannel
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

type messageFactoryPrivate struct {
	Author interface{} `json:"author"`
	//TODO: Add Member
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

func newMessage(state *clientState, data []byte) *Message {
	message := &Message{state: state}
	json.Unmarshal(data, message)
	private := new(messageFactoryPrivate)
	json.Unmarshal(data, private)
	userData, _ := json.Marshal(private.Author)
	message.Author = newUser(state, userData)
	if message.ChannelID != "" {
		message.Channel = state.GuildTextChannels[message.ChannelID]
	}
	if message.GuildID != "" {
		message.Guild = state.Guilds[message.GuildID]
	}
	mentionedGuildTextChannelsFactory(message)
	state.Messages[message.ID] = message
	return message
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
