package handler

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	cmsg "github.com/kellen-miller/gossip-gloomers/common/message"
	"github.com/kellen-miller/gossip-gloomers/common/node"
	"github.com/kellen-miller/gossip-gloomers/unique-ids/internal/message"
)

const (
	GenerateType      = "generate"
	GenerateTypeReply = GenerateType + "_ok"
)

type Generator struct {
	node *node.Node
}

func NewGenerator(n *node.Node) *Generator {
	uuid.EnableRandPool()
	return &Generator{
		node: n,
	}
}

func (g *Generator) Handle(ctx context.Context, msg *cmsg.Message) (*cmsg.Message, error) {
	body := new(message.GenerateBody)
	if err := json.Unmarshal(msg.Body, body); err != nil {
		return nil, err
	}

	body.ID = generateID()
	body.Type = GenerateTypeReply
	body.InReplyTo = body.MessageID

	bodyJ, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	return &cmsg.Message{
		Source:      g.node.ID,
		Destination: msg.Source,
		Body:        bodyJ,
	}, nil
}

func (g *Generator) Type() string {
	return GenerateType
}

func generateID() string {
	return uuid.New().String()
}
