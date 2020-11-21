package disgopher

//ClientState ...
type clientState struct {
	Guilds map[string]Guild
}

func (state *clientState) dispatch(name string, data map[string]interface{}) {
	switch name {
	case "GUILD_CREATE":
		newGuild(state, data)
	}
}
