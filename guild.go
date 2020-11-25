package disgopher

import (
	"encoding/json"
)

//Guild ...
type Guild struct {
	state        *clientState
	Description  string `json:"description"`
	MFALevel     int    `json:"mfa_level"`
	Region       string `json:"region"`
	AFKChannelID string `json:"afk_channel_id"`
	//VoiceStates []VoiceState `json:"voice_states"`
	//Presences []Presence `json:"presences"`
	AFKTimeout      int      `json:"afk_timeout"`
	MemberCount     int      `json:"member_count"`
	Icon            string   `json:"icon"`
	Splash          string   `json:"splash"`
	Features        []string `json:"features"`
	VaniryURLCode   string   `json:"vanity_url_code"`
	SystemChannelID string   `json:"system_channel_id"`
	//Roles []Role `json:"roles"`
	//Emojis []Emoji `json:"emojis"`
	JoinedAt                    string `json:"joined_at"`
	Name                        string `json:"name"`
	DefaultMessageNotifications int    `json:"default_message_notifications"`
	OwnerID                     string `json:"owner_id"`
	DiscoverySplash             string `json:"discovery_splash"`
	PremiumSubscriptionCount    int    `json:"premium_subscription_count"`
	ExplicitContentFilter       int    `json:"explicit_content_filter"`
	RulesChannelID              string `json:"rules_channel_id"`
	ApplicationID               string `json:"application_id"`
	MaxMembers                  int    `json:"max_members"`
	Unavailable                 bool   `json:"unavailable"`
	Large                       bool   `json:"large"`
	Lazy                        bool   `json:"lazy"`
	PublicUpdatesChannelID      string `json:"public_updates_channel_id"`
	PreferredLocale             string `json:"preferred_locale"`
	MaxVideoChannelUsers        int    `json:"max_video_channel_users"`
	PremiumTier                 int    `json:"premium_tier"`
	Banner                      string `json:"banner"`
	VerificationLevel           int    `json:"verification_level"`
	SystemChannelFlags          int    `json:"system_channel_flags"`
	ID                          string `json:"id"`
	TextChannels                map[string]*GuildTextChannel
	VoiceChannels               map[string]*GuildVoiceChannel
	Members                     map[string]*GuildMember
}

type guildFactoryPrivate struct {
	Channels []interface{} `json:"channels"`
	Members  []interface{} `json:"members"`
}

func channelFactory(guild *Guild, private *guildFactoryPrivate) {
	for index := range private.Channels {
		data, _ := json.Marshal(private.Channels[index])
		newBaseChannel(guild.state, data).upgrade(guild)
	}
}

func memberFactory(guild *Guild, private *guildFactoryPrivate) {
	for index := range private.Members {
		data, _ := json.Marshal(private.Members[index])
		newGuildMember(nil, guild, data)
	}
}

func newGuild(state *clientState, data []byte) *Guild {
	guild := &Guild{
		state:         state,
		TextChannels:  make(map[string]*GuildTextChannel),
		VoiceChannels: make(map[string]*GuildVoiceChannel),
		Members:       make(map[string]*GuildMember)}
	json.Unmarshal(data, guild)
	state.Guilds[guild.ID] = guild
	private := new(guildFactoryPrivate)
	json.Unmarshal(data, private)
	channelFactory(guild, private)
	memberFactory(guild, private)
	return guild
}

//GuildMember ...
type GuildMember struct {
	state *clientState
	Guild *Guild `json:"-"`
	User  *User  `json:"user"`
	Nick  string `json:"nick"`
	//Roles []
	JoinedAt     string `json:"joined_at"`
	PremiumSince string `json:"premium_since"`
	Deaf         bool   `json:"deaf"`
	Mute         bool   `json:"mute"`
}

func newGuildMember(user *User, guild *Guild, data []byte) *GuildMember {
	member := &GuildMember{state: guild.state, Guild: guild}
	json.Unmarshal(data, member)
	if user != nil {
		member.User = user
	} else if member.User != nil {
		member.User.state = guild.state
		guild.state.Users[member.User.ID] = member.User
	}
	guild.Members[member.User.ID] = member
	return member
}
