package message

const (
	InitType      = "init"
	InitReplyType = InitType + "_ok"
)

type Init struct {
	Base
	NodeID  string   `json:"node_id"`
	NodeIDs []string `json:"node_ids"`
}
