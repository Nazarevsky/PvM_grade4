package responses

// ReplaceRoundBracketsRequest Запит
type ReplaceRoundBracketsRequest struct {
	Message string `json:"message"`
}

// ReplaceRoundBracketsResponse Відповідь
type ReplaceRoundBracketsResponse struct {
	MessageProcessed string `json:"message_processed"`
}
