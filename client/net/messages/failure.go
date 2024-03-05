package messages

// Повідомлення про помилку
type FailureResponse struct {
	Message string `json:"result"`
}
