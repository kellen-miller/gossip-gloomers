package handler

import (
	"encoding/json"

	cmsg "github.com/kellen-miller/gossip-gloomers/common/message"
	"github.com/kellen-miller/gossip-gloomers/common/node"
)

const (
	TopologyType      = "topology"
	TopologyTypeReply = TopologyType + "_ok"
)

type Topology struct {
	n *node.Node
}

type TopologyBody struct {
	cmsg.BaseBody
	Topology map[string][]string `json:"topology,omitempty"`
}

func NewTopology(n *node.Node) *Topology {
	return &Topology{n: n}
}

func (t *Topology) Types() []string {
	return []string{TopologyType}
}

func (t *Topology) Handle(msg *cmsg.Message) (*cmsg.Message, error) {
	topologyBody := new(TopologyBody)
	if err := json.Unmarshal(msg.Body, topologyBody); err != nil {
		return nil, err
	}

	if neighbors, ok := topologyBody.Topology[t.n.ID]; ok {
		t.n.Neighbors = neighbors
	}

	replyBodyB, err := json.Marshal(&TopologyBody{
		BaseBody: cmsg.BaseBody{
			Type:      TopologyTypeReply,
			InReplyTo: topologyBody.MessageID,
		},
	})
	if err != nil {
		return nil, err
	}

	return &cmsg.Message{
		Source:      t.n.ID,
		Destination: msg.Source,
		Body:        replyBodyB,
	}, nil
}
