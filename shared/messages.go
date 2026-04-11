package shared

import "encoding/json"

type Packet struct {
	Type   string          `json:"type"`
	Action string          `json:"action"`
	Data   json.RawMessage `json:"data"`
}

type RoomState struct {
	OnlineUsers []string `json:"onlineUsers"`
}

type HeartBeat struct {
	HeartBeat string `json:"heartBeat"`
}

// --------------------------------------
// State Delta Messages
// --------------------------------------
type UserJoined struct {
	UserID string `json:"userID"`
}

type UserLeft struct {
	UserID string `json:"userID"`
}

type PostMsg struct {
	UserId string `json:"userID"`
	Msg    string `json:"msg"`
	Color  int    `json:"color"`
}
