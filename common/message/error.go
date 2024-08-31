package message

import (
	"encoding/json"
	"fmt"
	"os"
)

type ErrorBody struct {
	BaseBody
	Code int    `json:"code"`
	Text string `json:"text"`
}

func SendError(src string, dest string, inReplyTo int, code int, text string) {
	body := &ErrorBody{
		BaseBody: BaseBody{
			Type:      "error",
			InReplyTo: inReplyTo,
		},
		Code: code,
		Text: text,
	}

	bodyB, err := json.Marshal(body)
	if err != nil {
		fmt.Println("error marshalling error body:", err)
		return
	}

	msg := &Message{
		Source:      src,
		Destination: dest,
		Body:        bodyB,
	}

	msgB, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("error marshalling error message:", err)
		return
	}

	_, _ = fmt.Fprintf(os.Stderr, "%s\n", msgB)
}
