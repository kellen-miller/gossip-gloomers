package node

import (
	"encoding/json"

	"github.com/kellen-miller/gossip-gloomers/common/message"
)

const (
	InitType      = "init"
	InitReplyType = InitType + "_ok"
)

type Init struct {
	node *Node
}

type InitBody struct {
	message.BaseBody
	NodeID  string   `json:"node_id"`
	NodeIDs []string `json:"node_ids"`
}

func NewInit(n *Node) *Init {
	return &Init{
		node: n,
	}
}

func (i *Init) Types() []string {
	return []string{InitType}
}

func (i *Init) Handle(msg *message.Message) (*message.Message, error) {
	initBody := new(InitBody)
	if err := json.Unmarshal(msg.Body, initBody); err != nil {
		return nil, err
	}

	i.node.ID = initBody.NodeID
	i.node.IDs = initBody.NodeIDs

	replyBodyB, err := json.Marshal(&message.BaseBody{
		Type:      InitReplyType,
		MessageID: initBody.MessageID,
		InReplyTo: initBody.MessageID,
	})
	if err != nil {
		return nil, err
	}

	return &message.Message{
		Source:      i.node.ID,
		Destination: msg.Source,
		Body:        replyBodyB,
	}, nil
}
