package handler

import (
	"encoding/json"

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
	Message int `json:"message,omitempty"`
}

func NewBroadcast(n *node.Node, valsSeenChan chan int) *Broadcast {
	return &Broadcast{
		node:          n,
		nodesSeenChan: valsSeenChan,
	}
}

func (b *Broadcast) Type() string {
	return BroadcastType
}

func (b *Broadcast) Handle(msg *cmsg.Message) (*cmsg.Message, error) {
	broadcastBody := new(BroadcastBody)
	if err := json.Unmarshal(msg.Body, broadcastBody); err != nil {
		return nil, err
	}

	b.nodesSeenChan <- broadcastBody.Message

	replyBodyB, err := json.Marshal(&BroadcastBody{
		BaseBody: cmsg.BaseBody{
			Type: BroadcastReplyType,
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
