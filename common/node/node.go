package node

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/kellen-miller/gossip-gloomers/common/message"
	"github.com/ugurcsen/gods-generic/sets/hashset"
)

var ErrFireAndForget = errors.New("fire and forget")

type Handler interface {
	Handle(msg *message.Message) (*message.Message, error)
	Types() []string
}

type Node struct {
	ID                   string
	IDs                  []string
	Neighbors            []string
	typeToMessageIDsSeen map[string]*hashset.Set[int]
	handlers             map[string]Handler
}

func NewNode() *Node {
	n := &Node{
		typeToMessageIDsSeen: make(map[string]*hashset.Set[int]),
		handlers:             make(map[string]Handler),
	}

	n.RegisterHandlers(NewInit(n))
	return n
}

func (n *Node) RegisterHandlers(handlers ...Handler) {
	if n.handlers == nil {
		n.handlers = make(map[string]Handler, len(handlers))
	}

	for _, h := range handlers {
		for _, t := range h.Types() {
			n.handlers[t] = h
		}
	}
}

func (n *Node) Handle(msg *message.Message) (*message.Message, error) {
	baseBody := new(message.BaseBody)
	if err := json.Unmarshal(msg.Body, baseBody); err != nil {
		return nil, err
	}

	if n.checkMessageIDsSeen(baseBody.Type, baseBody.MessageID) {
		return nil, nil
	}

	handler, ok := n.handlers[baseBody.Type]
	if !ok {
		return nil, fmt.Errorf("no handler for message type %s", baseBody.Type)
	}

	return handler.Handle(msg)
}

func (n *Node) Listen(ctx context.Context) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan struct{})

	go func() {
		defer close(done)
		n.listen()
	}()

	select {
	case <-ctx.Done():
		println("Context cancelled. Shutting down...")
	case <-sigChan:
		println("Received interrupt signal. Shutting down...")
	case <-done:
		println("Message processing completed. Shutting down...")
	}
}

func (n *Node) checkMessageIDsSeen(kind string, msgID int) bool {
	seen, ok := n.typeToMessageIDsSeen[kind]
	if !ok {
		n.typeToMessageIDsSeen[kind] = hashset.New[int]()
		return false
	}

	return seen.Contains(msgID)
}

func (n *Node) listen() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		bytes := scanner.Bytes()

		msg := new(message.Message)
		if err := json.Unmarshal(bytes, msg); err != nil {
			_, _ = fmt.Fprintf(os.Stdout, "error unmarshalling message: %s\n", err)
			continue
		}

		resp, err := n.Handle(msg)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stdout, "error handling message: %s\n", err)
			continue
		}
		if resp == nil {
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
