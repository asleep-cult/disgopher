package disgopher

//Client ...
type Client struct {
	Token string
	Bot   bool
	Ws    *gateway
	HTTP  *HTTPSession
	State *clientState
}

//Start ...
func (client *Client) Start(token string) {
	client.HTTP.Token = token
	client.Ws.start(token)
}

//On ...
func (client *Client) On(name string, function interface{}) {
	client.State.registerEvent(name, function)
}

//NewClient ...
func NewClient() Client {
	state := &clientState{
		Guilds:             make(map[string]*Guild),
		GuildTextChannels:  make(map[string]*GuildTextChannel),
		GuildVoiceChannels: make(map[string]*GuildVoiceChannel),
		Messages:           make(map[string]*Message),
		Users:              make(map[string]*User),
		Events:             make(map[string][]interface{})}
	gateway := newGateway(state)
	rest := &HTTPSession{
		state:            state,
		ratelimitBuckets: make(map[string]*ratelimitBucket)}
	client := Client{Ws: gateway, State: state, HTTP: rest}
	return client
}
