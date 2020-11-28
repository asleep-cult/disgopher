package disgopher

import "sync"

//Client ...
type Client struct {
	Token string
	Bot   bool
	Ws    *gateway
	http  *HTTPSession
	State *clientState
}

//Start ...
func (client *Client) Start(token string) {
	client.http.Token = token
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
	http := &HTTPSession{
		state:                state,
		ratelimitBuckets:     make(map[string]*ratelimitBucket),
		globalRatelimitMutex: new(sync.Mutex)}
	state.Gateway = gateway
	state.HTTP = http
	client := Client{Ws: gateway, State: state, http: http}
	return client
}
