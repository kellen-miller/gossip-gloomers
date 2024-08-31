package handler

import (
	"context"
	"encoding/json"

	cmsg "github.com/kellen-miller/gossip-gloomers/common/message"
	cnode "github.com/kellen-miller/gossip-gloomers/common/node"
	"github.com/kellen-miller/gossip-gloomers/echo/internal/message"
)

const (
	EchoType      = "echo"
	EchoTypeReply = EchoType + "_ok"
)

type Echo struct {
	node *cnode.Node
}

func NewEcho(n *cnode.Node) *Echo {
	return &Echo{
		node: n,
	}
}

func (e *Echo) Type() string {
	return EchoType
}

func (e *Echo) Handle(_ context.Context, msg *cmsg.Message) (*cmsg.Message, error) {
	echoBody := new(message.EchoBody)
	if err := json.Unmarshal(msg.Body, echoBody); err != nil {
		return nil, err
	}

	echoBody.Type = EchoTypeReply
	echoBody.InReplyTo = echoBody.MessageID

	body, err := json.Marshal(echoBody)
	if err != nil {
		return nil, err
	}

	return &cmsg.Message{
		Source:      e.node.ID,
		Destination: msg.Source,
		Body:        body,
	}, nil
}
