package messages

import "encoding/json"

type Sum struct {
	Val1 int64 `json:"val1"`
	Val2 int64 `json:"val2"`
}

type SumResponse struct {
	Result int64 `json:"result"`
}

func NewSumMessageSerialized(val1, val2 int64) ([]byte, error) {
	failureResponseMessageBody, err := json.Marshal(Sum{Val1: val1, Val2: val2})
	if err != nil {
		return nil, err
	}

	message, err := json.Marshal(Message{Header: SumType, Body: failureResponseMessageBody})
	if err != nil {
		return nil, err
	}

	return message, nil
}
