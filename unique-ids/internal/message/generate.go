package message

import cmsg "github.com/kellen-miller/gossip-gloomers/common/message"

type GenerateBody struct {
	cmsg.BaseBody
	ID string `json:"id"`
}
