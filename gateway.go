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

type gateway struct {
	State             clientState
	ContinueBeating   bool
	HeartbeatInterval float64
	HearteatsSent     int
	HeartbeatsAcked   int
	LastSent          time.Time
	LastAcked         time.Time
	Sequence          int
	WebSocket         *websocket.Conn
	Token             string
}

type gatewayResponse struct {
	Name     interface{}            `json:"t"`
	Sequence interface{}            `json:"s"`
	Opcode   interface{}            `json:"op"`
	Data     map[string]interface{} `json:"d"`
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
		resp := new(gatewayResponse)
		json.Unmarshal(message, resp)
		if resp.Sequence != nil {
			gateway.Sequence = int(resp.Sequence.(float64))
		}
		if opcode, ok := resp.Opcode.(float64); ok {
			switch opcode {
			case opcodeDispatch:
				gateway.State.dispatch(resp.Name.(string), resp.Data)
			case opcodeHello:
				gateway.HeartbeatInterval = resp.Data["heartbeat_interval"].(float64)
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
	Interval := time.Duration(gateway.HeartbeatInterval) * time.Millisecond
	for gateway.ContinueBeating {
		gateway.sendHeartbeat()
		gateway.HearteatsSent++
		gateway.LastSent = time.Now()
		time.Sleep(Interval)
	}
}

func (gateway *gateway) sendIdentify() {
	payload := fmt.Sprintf("{\"op\":%v,\"d\":{\"token\":\"%s\",\"properties\":{\"$os\":\"Windows\",\"$browser\":\"Hinux\",\"$device\":\"Hinux\"}}}", opcodeIdentify, gateway.Token)
	err := gateway.WebSocket.WriteMessage(websocket.TextMessage, []byte(payload))
	if err != nil {
		panic(err)
	}
}

func (gateway *gateway) sendHeartbeat() {
	payload := fmt.Sprintf("{\"op\":%v,\"d\":null}", opcodeHeartbeat)
	err := gateway.WebSocket.WriteMessage(websocket.TextMessage, []byte(payload))
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
