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

//NewClient ...
func NewClient() Client {
	state := &clientState{Guilds: make(map[string]*Guild)}
	gateway := newGateway(state)
	client := Client{Ws: gateway, State: state}
	return client
}
