package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"server/net/messages"
)

// Створення та серіалізація повідомлення про суму двох чисел
func CreateSumResponse(sumMessage messages.Sum) ([]byte, error) {
	sumSerialized, err := json.Marshal(messages.SumResponse{Result: sumMessage.Val1 + sumMessage.Val2})
	if err != nil {
		return nil, errors.New(fmt.Sprint("Failed to serialize sum body response:", err))
	}

	serializedMessage, err := json.Marshal(messages.Message{Header: messages.SumResponseType, Body: sumSerialized})
	if err != nil {
		return nil, errors.New(fmt.Sprint("Failed to serialize sum response:", err))
	}

	return serializedMessage, nil
}
