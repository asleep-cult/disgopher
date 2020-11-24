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
	//Members []GuildMember `json:"members"`
}

//ToJSON ...
func (guild *Guild) ToJSON() ([]byte, error) {
	return json.Marshal(*guild)
}

func channelFactory(guild *Guild, data []byte) {
	var channels []interface{}
	channelStruct := struct {
		Channels []interface{} `json:"channels"`
	}{channels}
	json.Unmarshal(data, &channelStruct)
	channels = channelStruct.Channels
	for index := range channels {
		channel := ChannelBase{state: guild.state}
		channel.data, _ = json.Marshal(channels[index])
		switch newChannel := channel.upgrade(guild).(type) {
		case *GuildTextChannel:
			guild.TextChannels[newChannel.ID] = newChannel
		}
	}
}

func newGuild(state *clientState, data []byte) *Guild {
	guild := &Guild{state: state, TextChannels: make(map[string]*GuildTextChannel)}
	err := json.Unmarshal(data, guild)
	if err != nil {
		panic(err)
	}
	channelFactory(guild, data)
	return guild
}
