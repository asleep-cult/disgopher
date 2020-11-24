package disgopher

import "fmt"

//ClientState ...
type clientState struct {
	Guilds            map[string]*Guild
	GuildTextChannels map[string]*GuildTextChannel
	Events            map[string][]interface{}
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
		switch channel := newBaseChannel(state, data).upgrade().(type) {
		case *GuildTextChannel:
			fmt.Print(channel)
			if channel.Guild != nil {
				channel.Guild.TextChannels[channel.ID] = channel
			}
		case *GuildVoiceChannel:
			fmt.Print(channel)
			if channel.Guild != nil {
				channel.Guild.VoiceChannels[channel.ID] = channel
			}
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
	fmt.Print(state.Events[name])
}
