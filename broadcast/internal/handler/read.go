package handler

import (
	"context"
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
	n           *node.Node
	valsSeenSet *hashset.Set[int]
}

type ReadBody struct {
	cmsg.BaseBody
	Messages []int `json:"messages"`
}

func NewRead(ctx context.Context, n *node.Node, valsSeenChan chan int) *Read {
	r := &Read{
		n:           n,
		valsSeenSet: hashset.New[int](),
	}

	go r.processSeenChan(ctx, valsSeenChan)
	return r
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
		Messages: r.valsSeenSet.Values(),
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

func (r *Read) processSeenChan(ctx context.Context, seenChan chan int) {
	defer close(seenChan)

	for {
		select {
		case nodeID := <-seenChan:
			r.valsSeenSet.Add(nodeID)
		case <-ctx.Done():
			return
		}
	}
}
