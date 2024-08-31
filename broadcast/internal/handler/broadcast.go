package handler

import (
	"encoding/json"

	"github.com/kellen-miller/gossip-gloomers/broadcast/internal/message"
	cmsg "github.com/kellen-miller/gossip-gloomers/common/message"
	"github.com/kellen-miller/gossip-gloomers/common/node"
)

const (
	BroadcastType      = "broadcast"
	BroadcastReplyType = BroadcastType + "_ok"
)

type Broadcast struct {
	node          *node.Node
	nodesSeenChan chan int
}

type BroadcastBody struct {
	cmsg.BaseBody
	Message int `json:"message"`
}

func NewBroadcast(n *node.Node, nodesSeenChan chan int) *Broadcast {
	return &Broadcast{
		node:          n,
		nodesSeenChan: nodesSeenChan,
	}
}

func (b *Broadcast) Type() string {
	return BroadcastType
}

func (b *Broadcast) Handle(msg *cmsg.Message) (*cmsg.Message, error) {
	broadcastBody := new(message.BroadcastBody)
	if err := json.Unmarshal(msg.Body, broadcastBody); err != nil {
		return nil, err
	}

	b.nodesSeenChan <- broadcastBody.Message

	replyBody := &message.BroadcastBody{
		BaseBody: cmsg.BaseBody{
			Type: BroadcastReplyType,
		},
	}

	replyBodyB, err := json.Marshal(replyBody)
	if err != nil {
		return nil, err
	}

	return &cmsg.Message{
		Source:      b.node.ID,
		Destination: msg.Source,
		Body:        replyBodyB,
	}, nil
}
