package messages

type MessageType int

const (
	SumType MessageType = iota
	SumResponseType
	FailureType
)

// Загальна структура повідомлення
type Message struct {
	Header MessageType `json:"header"`
	Body   []byte      `json:"body"`
}
