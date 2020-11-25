package disgopher

//ClientState ...
type clientState struct {
	Guilds             map[string]*Guild
	GuildTextChannels  map[string]*GuildTextChannel
	GuildVoiceChannels map[string]*GuildVoiceChannel
	Messages           map[string]*Message
	Users              map[string]*User
	Events             map[string][]interface{}
}

func (state *clientState) dispatch(name string, data []byte) {
	switch name {
	case "GUILD_CREATE":
		guild := newGuild(state, data)
		for index := range state.Events["guild_create"] {
			event := GuildCreateEvent{Guild: guild}
			go state.Events["guild_create"][index].(func(GuildCreateEvent))(event)
		}
	case "CHANNEL_CREATE":
		newBaseChannel(state, data).upgrade()
	case "MESSAGE_CREATE":
		message := newMessage(state, data)
		for index := range state.Events["message_create"] {
			event := MessageCreateEvent{Message: message}
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
