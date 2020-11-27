package disgopher

//ClientState ...
type clientState struct {
	Guilds             map[string]*Guild
	GuildTextChannels  map[string]*GuildTextChannel
	GuildVoiceChannels map[string]*GuildVoiceChannel
	Messages           map[string]*Message
	Users              map[string]*User
	Events             map[string][]interface{}
	HTTP               *HTTPSession
	Gateway            *gateway
}

func (state *clientState) dispatch(name string, data []byte) {
	switch name {
	case "GUILD_CREATE":
		guild := newGuild(state, data)
		event := GuildCreateEvent{Guild: guild}
		for index := range state.Events["guild_create"] {
			go state.Events["guild_create"][index].(func(GuildCreateEvent))(event)
		}
	case "CHANNEL_CREATE":
		newBaseChannel(state, data).upgrade()
	case "MESSAGE_CREATE":
		message := newMessage(state, "", "", data)
		event := MessageCreateEvent{Message: message}
		for index := range state.Events["message_create"] {
			go state.Events["message_create"][index].(func(MessageCreateEvent))(event)
		}
	}
}

func (state *clientState) registerEvent(name string, function interface{}) {
	events := state.Events[name]
	if events == nil {
		var events []interface{}
		state.Events[name] = events
	}
	state.Events[name] = append(events, function)
}
