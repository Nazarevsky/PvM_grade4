package messages

const delim = '\n'

type MessageType int

const (
	SumType MessageType = iota
	SumResponseType
	FailureType
)

type Message struct {
	Header MessageType `json:"header"`
	Body   []byte      `json:"body"`
}

// Функція для додавання розділення для визначення кінця повідомлення
func AddDelim(data []byte) []byte {
	return append(data, delim)
}
