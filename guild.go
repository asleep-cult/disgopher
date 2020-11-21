package disgopher

//Guild ...
type Guild struct {
	Large                       bool
	PreferredLocale             string
	MaxVideoChannelUsers        int
	OwnerID                     string
	AFKTimeout                  int
	PremiumSubscriptionCount    int
	Name                        string
	Region                      string
	PremiumTier                 int
	MemberCount                 int
	PublicUpdatesChannelID      string
	DefaultMessageNotifications int
	SystemChannelFlags          int
	ExplicitContentFilter       int
	Lazy                        bool
	VerificationLevel           int
	SystemChannelID             string
	ID                          string
	Unavailable                 bool
	MFALevel                    int
}

//NewGuild ...
func newGuild(state *clientState, data map[string]interface{}) {
	guild := Guild{}
	large := data["large"]
	preferredLocale := data["preferred_locale"]
	maxVideoChannelUsers := data["max_video_channel_users"]
	//afkChannelID := data["afk_channel_id"]
	ownerID := data["owner_id"]
	afkTimeout := data["afk_timeout"]
	premiumSubscriptionCount := data["premium_subscription_count"]
	name := data["name"]
	region := data["region"]
	memberCount := data["member_count"]
	publicUpdatesChannelID := data["public_updates_channel_id"]
	defaultMessageNotifications := data["default_message_notifications"]
	systemChannelFlags := data["system_channel_flags"]
	explicitContentFilter := data["explicit_content_filter"]
	lazy := data["lazy"]
	//joinedAt := data["joined_at"]
	verificationLevel := data["verification_level"]
	systemChannelID := data["system_channel_id"]
	id := data["id"]
	unavailable := data["unavailable"]
	mfaLevel := data["mfa_level"]
	if large != nil {
		guild.Large = large.(bool)
	}
	if preferredLocale != nil {
		guild.PreferredLocale = preferredLocale.(string)
	}
	if maxVideoChannelUsers != nil {
		guild.MaxVideoChannelUsers = int(maxVideoChannelUsers.(float64))
	}
	if ownerID != nil {
		guild.OwnerID = ownerID.(string)
	}
	if afkTimeout != nil {
		guild.AFKTimeout = int(afkTimeout.(float64))
	}
	if premiumSubscriptionCount != nil {
		guild.PremiumSubscriptionCount = int(premiumSubscriptionCount.(float64))
	}
	if name != nil {
		guild.Name = name.(string)
	}
	if region != nil {
		guild.Region = region.(string)
	}
	if memberCount != nil {
		guild.MemberCount = int(memberCount.(float64))
	}
	if publicUpdatesChannelID != nil {
		guild.PublicUpdatesChannelID = publicUpdatesChannelID.(string)
	}
	if defaultMessageNotifications != nil {
		guild.DefaultMessageNotifications = int(defaultMessageNotifications.(float64))
	}
	if systemChannelFlags != nil {
		guild.SystemChannelFlags = int(systemChannelFlags.(float64))
	}
	if explicitContentFilter != nil {
		guild.ExplicitContentFilter = int(explicitContentFilter.(float64))
	}
	if lazy != nil {
		guild.Lazy = lazy.(bool)
	}
	if verificationLevel != nil {
		guild.VerificationLevel = int(verificationLevel.(float64))
	}
	if systemChannelID != nil {
		guild.SystemChannelID = systemChannelID.(string)
	}
	if id != nil {
		guild.ID = id.(string)
	}
	if unavailable != nil {
		guild.Unavailable = unavailable.(bool)
	}
	if mfaLevel != nil {
		guild.MFALevel = int(mfaLevel.(float64))
	}
	state.Guilds[guild.ID] = guild
}
