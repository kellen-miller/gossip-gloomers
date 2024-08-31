package message

import cmsg "github.com/kellen-miller/gossip-gloomers/common/message"

type EchoBody struct {
	cmsg.BaseBody
	Echo string `json:"echo"`
}
