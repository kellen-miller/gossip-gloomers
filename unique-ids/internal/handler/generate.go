package handler

import (
	"encoding/json"

	"github.com/google/uuid"
	cmsg "github.com/kellen-miller/gossip-gloomers/common/message"
	"github.com/kellen-miller/gossip-gloomers/common/node"
)

const (
	GenerateType      = "generate"
	GenerateReplyType = GenerateType + "_ok"
)

type Generator struct {
	node *node.Node
}

type GenerateBody struct {
	cmsg.BaseBody
	ID string `json:"id"`
}

func NewGenerator(n *node.Node) *Generator {
	uuid.EnableRandPool()
	return &Generator{
		node: n,
	}
}

func (g *Generator) Handle(msg *cmsg.Message) (*cmsg.Message, error) {
	genBody := new(GenerateBody)
	if err := json.Unmarshal(msg.Body, genBody); err != nil {
		return nil, err
	}

	replyBodyB, err := json.Marshal(&GenerateBody{
		BaseBody: cmsg.BaseBody{
			Type:      GenerateReplyType,
			MessageID: genBody.MessageID,
			InReplyTo: genBody.MessageID,
		},
		ID: generateID(),
	})
	if err != nil {
		return nil, err
	}

	return &cmsg.Message{
		Source:      g.node.ID,
		Destination: msg.Source,
		Body:        replyBodyB,
	}, nil
}

func (g *Generator) Types() []string {
	return []string{GenerateType}
}

func generateID() string {
	return uuid.New().String()
}
