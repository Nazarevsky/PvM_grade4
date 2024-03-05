package handler

import (
	"errors"
	"fmt"
	"server/net/messages"
)

// Створення та серіалізація повідомлення про помилку
func FormFailureMessage(messageText string) ([]byte, error) {
	serializedMessage, err := messages.NewFailureMessageSerialized(messageText)
	if err != nil {
		errText := fmt.Sprint("Failed to serialize failure message", err)
		return nil, errors.New(errText)
	}

	return serializedMessage, nil
}
