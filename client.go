package disgopher

//Client ...
type Client struct {
	Token string
	Bot   bool
	Ws    *gateway
	State *clientState
}

//Start ...
func (client *Client) Start(token string) {
	client.Ws.start(token)
}

//On ...
func (client *Client) On(name string, function interface{}) {
	client.State.registerEvent(name, function)
}

//NewClient ...
func NewClient() Client {
	state := &clientState{Guilds: make(map[string]*Guild), Events: make(map[string][]interface{})}
	gateway := newGateway(state)
	client := Client{Ws: gateway, State: state}
	return client
}
