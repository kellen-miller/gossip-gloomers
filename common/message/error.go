package message

type Error struct {
	Base
	Code int    `json:"code"`
	Text string `json:"text"`
}

func NewError(inReplyTo int, code int, text string) *Error {
	return &Error{
		Base: Base{
			Type:      "error",
			InReplyTo: inReplyTo,
		},
		Code: code,
		Text: text,
	}
}
