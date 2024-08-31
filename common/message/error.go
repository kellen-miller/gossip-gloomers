package message

type ErrorBody struct {
	BaseBody
	Code int    `json:"code"`
	Text string `json:"text"`
}

func NewError(inReplyTo int, code int, text string) *ErrorBody {
	return &ErrorBody{
		BaseBody: BaseBody{
			Type:      "error",
			InReplyTo: inReplyTo,
		},
		Code: code,
		Text: text,
	}
}
