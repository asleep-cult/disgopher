package main

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

const (
	gatewayURL = "wss://gateway.discord.gg?v=8"

	opcodeDispatch            = 0
	opcodeHeartbeat           = 1
	opcodeIdentify            = 2
	opcodePresenceUpdate      = 3
	opcodeVoiceStateUpdate    = 4
	opcodeResume              = 6
	opcodeReconnect           = 7
	opcodeRequestGuildMembers = 8
	opcodeInvalidSession      = 9
	opcodeHello               = 10
	opcodeHeartbeatAck        = 11
)

//Gateway ...
type Gateway struct {
	HeartbeatInterval float64
	Sequence          int
	WebSocket         *websocket.Conn
	Token             string
}

//GatewayResponse ...
type gatewayResponse struct {
	Name     interface{} `json:"t"`
	Sequence interface{} `json:"s"`
	Opcode   interface{} `json:"op"`
	Data     interface{} `json:"d"`
}

//connect ...
func (gateway *Gateway) connect(token string) {
	gateway.Token = token
	conn, _, err := websocket.DefaultDialer.Dial(gatewayURL, nil)
	gateway.WebSocket = conn
	if err != nil {
		fmt.Printf("FUCK ITS BROKEN")
	}
}

//describe ...
func describe(i interface{}) {
	fmt.Printf("(%v, %T)\n", i, i)
}

//listen ...
func (gateway *Gateway) listen() {
	for {
		_, message, err := gateway.WebSocket.ReadMessage()
		if err != nil {
			fmt.Printf("FUCK I BROKE IT")
		}
		resp := new(gatewayResponse)
		json.Unmarshal(message, resp)
		if resp.Sequence != nil {
			gateway.Sequence = resp.Sequence.(int)
		}
		describe(resp.Opcode)
		if resp.Opcode.(float64) == opcodeHello {
			gateway.HeartbeatInterval = resp.Data.(map[interface{}]interface{})["heartbeat_interval"].(float64)
			gateway.sendIdentify()
		}
	}
}

func (gateway *Gateway) sendIdentify() {
	payload := fmt.Sprintf("{\"op\":null,\"d\":{\"token\":%s,\"properties\": {\"$os\":\"Windows\",\"$browser\":\"Hinux\",\"$device\":\"Hinux\"}}", gateway.Token)
	fmt.Printf(payload)
}

func (gateway *Gateway) start() {
	gateway.connect("fhoiasdfhoiasdhf")
	gateway.listen()
}

func main() {
	gateway := Gateway{}
	gateway.start()
}
