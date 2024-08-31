package message

const (
	InitType      = "init"
	InitReplyType = InitType + "_ok"
)

type InitBody struct {
	BaseBody
	NodeID  string   `json:"node_id"`
	NodeIDs []string `json:"node_ids"`
}
