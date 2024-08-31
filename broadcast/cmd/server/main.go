package main

import (
	"context"

	"github.com/kellen-miller/gossip-gloomers/common/node"
)

func main() {
	n := node.NewNode()
	n.Start(context.Background())
}
