package net

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"server/net/handler"
	"server/net/messages"
)

// Обробник TCP підключень
func handleTCPConnection(conn net.Conn) {
	reader := bufio.NewReader(conn)

	for {
		message, err := reader.ReadBytes('\n')
		if err != nil {
			return
		}

		serializedResponse := handleMessageBytes(message)
		_, err = conn.Write(serializedResponse)
		if err != nil {
			println(fmt.Sprint("failed to write a tcp response", err))
		}
	}
}

// Обробник UDP підключень
func handleUDPConnection(message []byte, conn *net.UDPConn, addr *net.UDPAddr) {
	serializedResponse := handleMessageBytes(message)
	_, err := conn.WriteToUDP(serializedResponse, addr)
	if err != nil {
		println(fmt.Sprint("failed to write a udp response: ", err))
	}
}

// Обробник повідомлень (серіалізований вигляд)
func handleMessageBytes(message []byte) []byte {
	var parsedMessage messages.Message
	err := handleNewMessage(message, &parsedMessage)
	if err != nil {
		errText := fmt.Sprint("failed to handle message bytes: ", err)
		println(errText)

		// Повідомлення клієнта про помилку, якщо невдається обробити нове повідомлення
		failureSerialized, failureErr := handler.FormFailureMessage(errText)
		if failureErr != nil {
			println(failureErr.Error())
			return nil
		}
		return failureSerialized
	}

	// Обробка різних типів повідомлень
	switch parsedMessage.Header {
	case messages.SumType:
		var sumMessage messages.Sum
		err = handleNewMessage(parsedMessage.Body, &sumMessage)
		if err != nil {
			errText := fmt.Sprint("failed to handle message bytes: ", err)
			println(errText)

			failureSerialized, failureErr := handler.FormFailureMessage(errText)
			if failureErr != nil {
				println(failureErr.Error())
				return nil
			}
			return failureSerialized
		}

		// Створення відповіді на повідомлення з запитом розрахунку суми
		sumResponse, err := handler.CreateSumResponse(sumMessage)
		if err != nil {
			errText := fmt.Sprint("failed to create sum response", err)
			println(errText)

			failureSerialized, failureErr := handler.FormFailureMessage(errText)
			if failureErr != nil {
				println(failureErr.Error())
				return nil
			}
			return failureSerialized
		}
		return sumResponse
	default:
		// Створення повідомлення про помилку у випадку, якщо отримана команда невідома
		errText := fmt.Sprint("unknown command")
		println(errText)

		failureSerialized, failureErr := handler.FormFailureMessage(errText)
		if failureErr != nil {
			println(failureErr.Error())
			return nil
		}
		return failureSerialized
	}
}

// Обробка та десеріалізація повідомлень
func handleNewMessage(bytes []byte, message any) error {
	err := json.Unmarshal(bytes, &message)
	if err != nil {
		return errors.New(fmt.Sprint("failed to handle new message: ", err))
	}

	return nil
}

// Запуск слухача відповідно до протоколу
func Run(port int, protocol string) {
	address := fmt.Sprintf("localhost:%d", port)

	if protocol == "tcp" {
		startTcp(address)
	} else if protocol == "udp" {
		startUdp(address)
	} else {
		println("Unknown protocol, aborting")
	}
}

// Запуск TCP слухача
func startTcp(address string) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		fmt.Println("Error resolving TCP address:", err)
		return
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer listener.Close()
	fmt.Printf("Server listening on %s, tcp \n", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		// Обробка TCP підключення
		go handleTCPConnection(conn)
	}
}

// Запуск UDP слухача
func startUdp(address string) {
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		fmt.Println("Error resolving TCP address:", err)
		return
	}

	connection, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer connection.Close()
	fmt.Printf("Server listening on %s, udp \n", address)

	for {
		// Максимальний розмір повідомлення - 1024 байти
		message := make([]byte, 1024)
		size, addr, err := connection.ReadFromUDP(message)
		if err != nil {
			fmt.Println("Error reading from UDP connection:", err)
			continue
		}

		// Обробка UDP підключення
		go handleUDPConnection(message[:size], connection, addr)
	}
}
