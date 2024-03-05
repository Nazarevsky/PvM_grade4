package controller

import (
	"client/net/messages"
	"encoding/json"
	"errors"
	"fmt"
)

// Обробник відповіді
func HandleResponse(data []byte, assumedMessage any) (messages.FailureResponse, error) {
	var receivedMessage messages.Message
	err := json.Unmarshal(data, &receivedMessage)
	if err != nil {
		return messages.FailureResponse{}, errors.New(fmt.Sprintf("failed to parce received message: %s", err))
	}

	// В залежності від хедеру визначається спосіб для обробки повідомлення
	switch receivedMessage.Header {
	case messages.SumResponseType:
		sumMessage, ok := assumedMessage.(*messages.SumResponse)
		if !ok {
			return messages.FailureResponse{}, errors.New(fmt.Sprintf("failed to create a pointer to SumResponse: %s", err))
		}

		err = json.Unmarshal(receivedMessage.Body, &sumMessage)
		if err != nil {
			return messages.FailureResponse{}, errors.New(fmt.Sprintf("failed to parse sum body: %s", err))
		}

		return messages.FailureResponse{}, nil
	case messages.FailureType:
		var failureMessage messages.FailureResponse
		err = json.Unmarshal(receivedMessage.Body, &failureMessage)
		if err != nil {
			return messages.FailureResponse{}, errors.New(fmt.Sprintf("failed to parse failure body: %s", err))
		}

		return failureMessage, nil
	default:
		return messages.FailureResponse{}, errors.New("failed to handle unknown message")
	}
}
