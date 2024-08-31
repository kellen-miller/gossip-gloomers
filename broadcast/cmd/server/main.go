package main

import (
	"context"

	"github.com/kellen-miller/gossip-gloomers/broadcast/internal/handler"
	"github.com/kellen-miller/gossip-gloomers/common/node"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		n            = node.NewNode()
		valsSeenChan = make(chan int)
	)

	n.RegisterHandlers(
		handler.NewBroadcast(n, valsSeenChan),
		handler.NewRead(ctx, n, valsSeenChan),
		handler.NewTopology(n),
	)

	n.Listen(ctx)
}
