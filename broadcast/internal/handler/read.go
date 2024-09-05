package handler

import (
	"encoding/json"

	cmsg "github.com/kellen-miller/gossip-gloomers/common/message"
	"github.com/kellen-miller/gossip-gloomers/common/node"
	"github.com/ugurcsen/gods-generic/sets/hashset"
)

const (
	ReadType      = "read"
	ReadReplyType = ReadType + "_ok"
)

type Read struct {
	n            *node.Node
	messagesSeen *hashset.Set[int]
}

type ReadBody struct {
	cmsg.BaseBody
	Messages []int `json:"messages"`
}

func NewRead(n *node.Node, messagesSeen *hashset.Set[int]) *Read {
	return &Read{
		n:            n,
		messagesSeen: messagesSeen,
	}
}

func (r *Read) Types() []string {
	return []string{ReadType}
}

func (r *Read) Handle(msg *cmsg.Message) (*cmsg.Message, error) {
	readBody := new(ReadBody)
	if err := json.Unmarshal(msg.Body, readBody); err != nil {
		return nil, err
	}

	replyBodyB, err := json.Marshal(&ReadBody{
		BaseBody: cmsg.BaseBody{
			Type:      ReadReplyType,
			InReplyTo: readBody.MessageID,
		},
		Messages: r.messagesSeen.Values(),
	})
	if err != nil {
		return nil, err
	}

	return &cmsg.Message{
		Source:      r.n.ID,
		Destination: msg.Source,
		Body:        replyBodyB,
	}, nil
}
