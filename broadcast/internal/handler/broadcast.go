package handler

import (
	"encoding/json"
	"fmt"
	"os"

	cmsg "github.com/kellen-miller/gossip-gloomers/common/message"
	"github.com/kellen-miller/gossip-gloomers/common/node"
)

const (
	BroadcastType      = "broadcast"
	BroadcastReplyType = BroadcastType + "_ok"
)

type Broadcast struct {
	node         *node.Node
	messagesChan chan int
}

type BroadcastBody struct {
	cmsg.BaseBody
	Message int `json:"message,omitempty"`
}

func NewBroadcast(n *node.Node, msgsChan chan int) *Broadcast {
	return &Broadcast{
		node:         n,
		messagesChan: msgsChan,
	}
}

func (b *Broadcast) Types() []string {
	return []string{BroadcastType, BroadcastReplyType}
}

func (b *Broadcast) Handle(msg *cmsg.Message) (*cmsg.Message, error) {
	broadcastBody := new(BroadcastBody)
	if err := json.Unmarshal(msg.Body, broadcastBody); err != nil {
		return nil, err
	}

	if broadcastBody.Type == BroadcastReplyType {
		return nil, nil
	}

	b.messagesChan <- broadcastBody.Message

	for _, neighbor := range b.node.Neighbors {
		neighborB, err := json.Marshal(&BroadcastBody{
			BaseBody: cmsg.BaseBody{
				Type: BroadcastType,
			},
			Message: broadcastBody.Message,
		})
		if err != nil {
			return nil, err
		}

		msgB, err := json.Marshal(&cmsg.Message{
			Source:      b.node.ID,
			Destination: neighbor,
			Body:        neighborB,
		})
		if err != nil {
			return nil, err
		}

		_, _ = fmt.Fprintf(os.Stdout, "%s\n", msgB)
	}

	replyBodyB, err := json.Marshal(&BroadcastBody{
		BaseBody: cmsg.BaseBody{
			Type:      BroadcastReplyType,
			MessageID: broadcastBody.MessageID,
			InReplyTo: broadcastBody.MessageID,
		},
	})
	if err != nil {
		return nil, err
	}

	return &cmsg.Message{
		Source:      b.node.ID,
		Destination: msg.Source,
		Body:        replyBodyB,
	}, nil
}
