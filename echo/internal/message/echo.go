package message

import cmsg "github.com/kellen-miller/gossip-gloomers/common/message"

type Echo struct {
	cmsg.Base
	Echo string `json:"echo"`
}
