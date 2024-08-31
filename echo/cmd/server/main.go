package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"

	chandler "github.com/kellen-miller/gossip-gloomers/common/handler"
	cmsg "github.com/kellen-miller/gossip-gloomers/common/message"
	cnode "github.com/kellen-miller/gossip-gloomers/common/node"
	"github.com/kellen-miller/gossip-gloomers/echo/internal/handler"
)

func main() {
	n := cnode.NewNode()

	n.RegisterHandlers(
		chandler.NewInit(n),
		handler.NewEcho(n),
	)

	ctx := context.Background()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		bytes := scanner.Bytes()

		msg := new(cmsg.Message)
		if err := json.Unmarshal(bytes, msg); err != nil {
			_, _ = fmt.Fprintf(os.Stdout, "error unmarshalling message: %s\n", err)
			continue
		}

		resp, err := n.Handle(ctx, msg)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stdout, "error handling message: %s\n", err)
			continue
		}

		respB, err := json.Marshal(resp)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stdout, "error marshalling response: %s\n", err)
			continue
		}

		_, _ = fmt.Fprintf(os.Stdout, "%s\n", respB)
	}

	if err := scanner.Err(); err != nil {
		_, _ = fmt.Fprintf(os.Stdout, "error reading input: %s\n", err)
	}
}
