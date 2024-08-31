package handler

import (
	"context"

	cmsg "github.com/kellen-miller/gossip-gloomers/common/message"
	"github.com/kellen-miller/gossip-gloomers/common/node"
	"github.com/ugurcsen/gods-generic/sets/hashset"
)

const (
	ReadType       = "read"
	ReadeReplyType = ReadType + "_ok"
)

type Read struct {
	n            *node.Node
	nodesSeenSet *hashset.Set[int]
}

type ReadBody struct {
	cmsg.BaseBody
}

func NewRead(ctx context.Context, n *node.Node, nodesSeenChan chan int) *Read {
	r := &Read{
		n:            n,
		nodesSeenSet: hashset.New[int](),
	}

	go r.processSeenChan(ctx, nodesSeenChan)
	return r
}

func (r *Read) Type() string {
	return ReadType
}

func (r *Read) Handle(msg *cmsg.Message) (*cmsg.Message, error) {
	return nil, nil
}

func (r *Read) processSeenChan(ctx context.Context, seenChan chan int) {
	defer close(seenChan)

	for {
		select {
		case nodeID := <-seenChan:
			r.nodesSeenSet.Add(nodeID)
		case <-ctx.Done():
			return
		}
	}
}
