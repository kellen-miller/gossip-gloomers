package node

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/kellen-miller/gossip-gloomers/common/message"
)

type Handler interface {
	Handle(msg *message.Message) (*message.Message, error)
	Type() string
}

type Node struct {
	ID       string
	IDs      []string
	handlers map[string]Handler
}

func NewNode() *Node {
	n := &Node{handlers: make(map[string]Handler)}
	n.RegisterHandlers(NewInit(n))
	return n
}

func (n *Node) RegisterHandlers(handlers ...Handler) {
	if n.handlers == nil {
		n.handlers = make(map[string]Handler, len(handlers))
	}

	for _, h := range handlers {
		n.handlers[h.Type()] = h
	}
}

func (n *Node) Handle(ctx context.Context, msg *message.Message) (*message.Message, error) {
	base := new(message.BaseBody)
	if err := json.Unmarshal(msg.Body, base); err != nil {
		return nil, err
	}

	handler, ok := n.handlers[base.Type]
	if !ok {
		return nil, fmt.Errorf("no handler for message type %s", base.Type)
	}

	return handler.Handle(msg)
}

func (n *Node) Start(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan struct{})

	go func() {
		defer close(done)
		n.listen(ctx)
	}()

	select {
	case <-ctx.Done():
		fmt.Println("Context cancelled. Shutting down...")
	case <-sigChan:
		fmt.Println("Received interrupt signal. Shutting down...")
		cancel()
	case <-done:
		fmt.Println("Message processing completed. Shutting down...")
	}

	<-done
}

func (n *Node) listen(ctx context.Context) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		bytes := scanner.Bytes()

		msg := new(message.Message)
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
