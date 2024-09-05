package handler

import (
	"encoding/json"
	"fmt"
	"os"

	cmsg "github.com/kellen-miller/gossip-gloomers/common/message"
	"github.com/kellen-miller/gossip-gloomers/common/node"
	"github.com/ugurcsen/gods-generic/sets/hashset"
)

const (
	BroadcastType      = "broadcast"
	BroadcastReplyType = BroadcastType + "_ok"
)

type Broadcast struct {
	node         *node.Node
	messagesSeen *hashset.Set[int]
}

type BroadcastBody struct {
	cmsg.BaseBody
	Message int `json:"message,omitempty"`
}

func NewBroadcast(n *node.Node, messagesSeen *hashset.Set[int]) *Broadcast {
	return &Broadcast{
		node:         n,
		messagesSeen: messagesSeen,
	}
}

func (b *Broadcast) Types() []string {
	return []string{BroadcastType, BroadcastReplyType}
}

func (b *Broadcast) Handle(msg *cmsg.Message) (*cmsg.Message, error) {
	broadcastBody := new(BroadcastBody)
	if err := json.Unmarshal(msg.Body, broadcastBody); err != nil {
		return nil, err
	}

	if !b.messagesSeen.Contains(broadcastBody.Message) {
		b.messagesSeen.Add(broadcastBody.Message)

		if err := b.notifyNeighbors(broadcastBody); err != nil {
			return nil, err
		}
	}

	if broadcastBody.MessageID == 0 || broadcastBody.Type == BroadcastReplyType {
		return nil, nil
	}

	replyBodyB, err := json.Marshal(&BroadcastBody{
		BaseBody: cmsg.BaseBody{
			Type:      BroadcastReplyType,
			MessageID: broadcastBody.MessageID,
			InReplyTo: broadcastBody.MessageID,
		},
	})
	if err != nil {
		return nil, err
	}

	return &cmsg.Message{
		Source:      b.node.ID,
		Destination: msg.Source,
		Body:        replyBodyB,
	}, nil
}

func (b *Broadcast) notifyNeighbors(bb *BroadcastBody) error {
	for _, neighbor := range b.node.Neighbors {
		neighborB, err := json.Marshal(&BroadcastBody{
			BaseBody: cmsg.BaseBody{
				Type: BroadcastType,
			},
			Message: bb.Message,
		})
		if err != nil {
			return err
		}

		msgB, err := json.Marshal(&cmsg.Message{
			Source:      b.node.ID,
			Destination: neighbor,
			Body:        neighborB,
		})
		if err != nil {
			return err
		}

		_, _ = fmt.Fprintf(os.Stdout, "%s\n", msgB)
	}

	return nil
}
