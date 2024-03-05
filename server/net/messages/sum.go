package messages

// Запит на отримання суми двох чисел
type Sum struct {
	Val1 int64 `json:"val1"`
	Val2 int64 `json:"val2"`
}

// Відповідь - сума двох чисел
type SumResponse struct {
	Result int64 `json:"result"`
}
