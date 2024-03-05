package controller

import (
	"client/net/messages"
	"errors"
	"fmt"
	"net"
	"strings"
)

type MessageController struct {
	Protocol string
	Address  string
}

func New(protocol, address string) (MessageController, error) {
	protocolLower := strings.ToLower(protocol)
	if protocolLower != "tcp" && protocolLower != "udp" {
		return MessageController{}, errors.New("no such a protocol")
	}

	return MessageController{
		Protocol: protocol,
		Address:  address,
	}, nil
}

// Відправка повідомлення про суму
func (mc *MessageController) SendSum(num1, num2 int64) ([]byte, error) {
	conn, err := net.Dial(mc.Protocol, mc.Address)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to dial %s: %s", mc.Address, err))
	}
	defer conn.Close()

	// Створення нового серіалізованого повідомлення про суму двох чисел
	sumSerialized, err := messages.NewSumMessageSerialized(num1, num2)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to create serialized sum message: %s", err))
	}

	_, err = conn.Write(messages.AddDelim(sumSerialized))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to send sum message: %s", err))
	}

	// Встановлення максимального розміру каналу
	serializedMessage := make([]byte, 1024)
	size, err := conn.Read(serializedMessage)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to receieve a response: %s", err))
	}

	return serializedMessage[:size], nil
}
