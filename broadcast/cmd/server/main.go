package main

import (
	"context"

	"github.com/kellen-miller/gossip-gloomers/broadcast/internal/handler"
	"github.com/kellen-miller/gossip-gloomers/common/node"
	"github.com/ugurcsen/gods-generic/sets/hashset"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		n            = node.NewNode()
		messagesSeen = hashset.New[int]()
	)
	n.RegisterHandlers(
		handler.NewBroadcast(n, messagesSeen),
		handler.NewRead(n, messagesSeen),
		handler.NewTopology(n),
	)

	n.Listen(ctx)
}
