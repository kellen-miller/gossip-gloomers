package main

import (
	"context"

	"github.com/kellen-miller/gossip-gloomers/common/node"
	"github.com/kellen-miller/gossip-gloomers/echo/internal/handler"
)

func main() {
	n := node.NewNode()
	n.RegisterHandlers(handler.NewEcho(n))
	n.Start(context.Background())
}
