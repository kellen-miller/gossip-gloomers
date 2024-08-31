package message

import "encoding/json"

type Message struct {
	Source      string          `json:"src"`
	Destination string          `json:"dest"`
	Body        json.RawMessage `json:"body"`
}

type Base struct {
	Type      string `json:"type"`
	MessageID int    `json:"msg_id,omitempty"`
	InReplyTo int    `json:"in_reply_to,omitempty"`
}
