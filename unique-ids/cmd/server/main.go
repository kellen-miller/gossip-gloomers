package main

import (
	"context"

	"github.com/kellen-miller/gossip-gloomers/common/node"
	"github.com/kellen-miller/gossip-gloomers/unique-ids/internal/handler"
)

func main() {
	n := node.NewNode()
	n.RegisterHandlers(handler.NewGenerator(n))
	n.Start(context.Background())
}
