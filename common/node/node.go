package node

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kellen-miller/gossip-gloomers/common/message"
)

type Handler interface {
	Handle(ctx context.Context, msg *message.Message) (*message.Message, error)
	Type() string
}

type Node struct {
	ID       string
	IDs      []string
	handlers map[string]Handler
}

func NewNode() *Node {
	return &Node{
		handlers: make(map[string]Handler),
	}
}

func (n *Node) RegisterHandlers(handler ...Handler) {
	for _, h := range handler {
		n.handlers[h.Type()] = h
	}
}

func (n *Node) Handle(ctx context.Context, msg *message.Message) (*message.Message, error) {
	base := new(message.Base)
	if err := json.Unmarshal(msg.Body, base); err != nil {
		return nil, err
	}

	handler, ok := n.handlers[base.Type]
	if !ok {
		return nil, fmt.Errorf("no handler for message type %s", base.Type)
	}

	return handler.Handle(ctx, msg)
}
