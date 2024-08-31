package handler

import (
	"encoding/json"

	cmsg "github.com/kellen-miller/gossip-gloomers/common/message"
	cnode "github.com/kellen-miller/gossip-gloomers/common/node"
)

const (
	EchoType      = "echo"
	EchoReplyType = EchoType + "_ok"
)

type Echo struct {
	node *cnode.Node
}

type EchoBody struct {
	cmsg.BaseBody
	Echo string `json:"echo"`
}

func NewEcho(n *cnode.Node) *Echo {
	return &Echo{
		node: n,
	}
}

func (e *Echo) Type() string {
	return EchoType
}

func (e *Echo) Handle(msg *cmsg.Message) (*cmsg.Message, error) {
	echoBody := new(EchoBody)
	if err := json.Unmarshal(msg.Body, echoBody); err != nil {
		return nil, err
	}

	echoBodyB, err := json.Marshal(&EchoBody{
		BaseBody: cmsg.BaseBody{
			Type:      EchoReplyType,
			InReplyTo: echoBody.MessageID,
		},
		Echo: echoBody.Echo,
	})
	if err != nil {
		return nil, err
	}

	return &cmsg.Message{
		Source:      e.node.ID,
		Destination: msg.Source,
		Body:        echoBodyB,
	}, nil
}
