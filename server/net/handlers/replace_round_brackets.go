package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	responses "server/net/data"
	"strings"
)

func ReplaceRoundBrackets(w http.ResponseWriter, r *http.Request) {
	// Парсинг запиту
	var req responses.ReplaceRoundBracketsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Println("Error:", err)
		http.Error(w, "Error getting a request", http.StatusBadRequest)
		return
	}

	// Перевірка запиту на те, що він порожній
	if strings.ReplaceAll(req.Message, " ", "") == "" {
		fmt.Println("Empty message passed")
		http.Error(w, "The message is empty", http.StatusBadRequest)
		return
	}

	// Формування відповіді з повідомленням, де круглі дужки замінено на квадратні
	var response = responses.ReplaceRoundBracketsResponse{
		MessageProcessed: strings.ReplaceAll(strings.ReplaceAll(req.Message, "(", "["), ")", "]"),
	}

	// Серіалізація відповіді у json структуру
	res, err := json.Marshal(response)
	if err != nil {
		fmt.Println("Error:", err)
		http.Error(w, "Error marshalling a response", http.StatusInternalServerError)
		return
	}

	// Повернення відповіді
	_, err = w.Write(res)
	if err != nil {
		fmt.Println("Error:", err)
		http.Error(w, "Error marshalling a response", http.StatusBadGateway)
		return
	}
}
