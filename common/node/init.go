package node

import (
	"context"
	"encoding/json"

	"github.com/kellen-miller/gossip-gloomers/common/message"
)

type Init struct {
	node *Node
}

func NewInit(n *Node) *Init {
	return &Init{
		node: n,
	}
}

func (i *Init) Type() string {
	return message.InitType
}

func (i *Init) Handle(_ context.Context, msg *message.Message) (*message.Message, error) {
	initMsg := new(message.InitBody)
	if err := json.Unmarshal(msg.Body, initMsg); err != nil {
		return nil, err
	}

	i.node.ID = initMsg.NodeID
	i.node.IDs = initMsg.NodeIDs

	bodyB, err := json.Marshal(&message.BaseBody{
		Type:      message.InitReplyType,
		InReplyTo: initMsg.MessageID,
	})
	if err != nil {
		return nil, err
	}

	return &message.Message{
		Source:      i.node.ID,
		Destination: msg.Source,
		Body:        bodyB,
	}, nil
}
