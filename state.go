package disgopher

//ClientState ...
type clientState struct {
	Guilds            map[string]*Guild
	GuildTextChannels map[string]*GuildTextChannel
}

func (state *clientState) dispatch(name string, data []byte) {
	switch name {
	case "GUILD_CREATE":
		newGuild(state, data)
	}
}
