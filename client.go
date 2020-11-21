package disgopher

//Client ...
type Client struct {
	Token string
	Bot   bool
	Ws    gateway
	State clientState
}

//Start ...
func (client *Client) Start(token string) {
	client.Ws.start(token)
}

//NewClient ...
func NewClient() Client {
	gateway := gateway{}
	client := Client{Ws: gateway}
	client.State.Guilds = make(map[string]Guild)
	gateway.State = client.State
	return client
}
