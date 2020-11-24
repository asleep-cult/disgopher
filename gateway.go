package disgopher

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

const (
	gatewayURL = "wss://gateway.discord.gg"

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

type heartbeatRequest struct {
	Opcode int         `json:"op"`
	Data   interface{} `json:"d"`
}

type identifyProperties struct {
	OS      string `json:"$os"`
	Browser string `json:"$browser"`
	Device  string `json:"$device"`
}

type identifyData struct {
	Token      string             `json:"token"`
	Properties identifyProperties `json:"properties"`
}

type identifyRequest struct {
	Opcode int          `json:"op"`
	Data   identifyData `json:"d"`
}

type baseResponse struct {
	Name     interface{} `json:"t"`
	Sequence interface{} `json:"s"`
	Opcode   interface{} `json:"op"`
	Data     interface{} `json:"d"`
}

type helloResponse struct {
	HelloData struct {
		HeartbeatInterval float64 `json:"heartbeat_interval"`
	} `json:"d"`
}

type gateway struct {
	ContinueBeating bool

	HearteatsSent   int
	HeartbeatsAcked int
	Sequence        int

	Token string

	HeartbeatPayload []byte

	HeartbeatInterval time.Duration
	LastSent          time.Time
	LastAcked         time.Time

	State     *clientState
	WebSocket *websocket.Conn
}

func (gateway *gateway) connect(token string) {
	gateway.Token = token
	conn, _, err := websocket.DefaultDialer.Dial(gatewayURL, nil)
	gateway.WebSocket = conn
	if err != nil {
		panic(err)
	}
}

func (gateway *gateway) listen() {
	for {
		_, message, err := gateway.WebSocket.ReadMessage()
		if err != nil {
			panic(err)
		}
		resp := new(baseResponse)
		json.Unmarshal(message, resp)
		if resp.Sequence != nil {
			gateway.Sequence = int(resp.Sequence.(float64))
		}
		if opcode, ok := resp.Opcode.(float64); ok {
			switch opcode {
			case opcodeDispatch:
				data, _ := json.Marshal(resp.Data)
				gateway.State.dispatch(resp.Name.(string), data)
			case opcodeHello:
				helloResp := new(helloResponse)
				json.Unmarshal(message, helloResp)
				gateway.HeartbeatInterval = time.Duration(helloResp.HelloData.HeartbeatInterval) * time.Millisecond
				gateway.sendIdentify()
				if !gateway.ContinueBeating {
					go gateway.startBeating()
				}
			case opcodeHeartbeat:
				gateway.sendHeartbeat()
			case opcodeHeartbeatAck:
				gateway.HeartbeatsAcked++
				gateway.LastAcked = time.Now()
			default:
				fmt.Printf("unknown opcode? (%v %T)", opcode, opcode)
			}
		}
	}
}

func (gateway *gateway) startBeating() {
	gateway.ContinueBeating = true
	for gateway.ContinueBeating {
		gateway.sendHeartbeat()
		gateway.HearteatsSent++
		gateway.LastSent = time.Now()
		time.Sleep(gateway.HeartbeatInterval)
	}
}

func (gateway *gateway) sendIdentify() {
	properties := identifyProperties{OS: "Windows", Browser: "Hinux", Device: "Hinux"}
	data := identifyData{Token: gateway.Token, Properties: properties}
	identifyReq := identifyRequest{Opcode: opcodeIdentify, Data: data}
	payload, jsonerr := json.Marshal(identifyReq)
	if jsonerr != nil {
		panic(jsonerr)
	}
	err := gateway.WebSocket.WriteMessage(websocket.TextMessage, []byte(payload))
	if err != nil {
		panic(err)
	}
}

func (gateway *gateway) sendHeartbeat() {
	err := gateway.WebSocket.WriteMessage(websocket.TextMessage, gateway.HeartbeatPayload)
	if err != nil {
		panic(err)
	}
}

func (gateway *gateway) start(token string) {
	gateway.connect(token)
	gateway.listen()
}

//Latency ...
func (gateway *gateway) Latency() time.Duration {
	return gateway.LastAcked.Sub(gateway.LastSent)
}

func newGateway(state *clientState) *gateway {
	data, _ := json.Marshal(heartbeatRequest{Opcode: opcodeHeartbeat, Data: nil})
	gateway := gateway{State: state, HeartbeatPayload: data}
	return &gateway
}
