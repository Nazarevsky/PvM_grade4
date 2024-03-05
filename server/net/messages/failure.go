package messages

import "encoding/json"

// Повідомлення про помилку
type FailureResponse struct {
	Message string `json:"result"`
}

// Створення нового серіалізованого повідомлення про помилку
func NewFailureMessageSerialized(failureText string) ([]byte, error) {
	failureResponseMessageBody, err := json.Marshal(FailureResponse{Message: failureText})
	if err != nil {
		return nil, err
	}

	message, err := json.Marshal(Message{Header: FailureType, Body: failureResponseMessageBody})
	if err != nil {
		return nil, err
	}

	return message, nil
}
