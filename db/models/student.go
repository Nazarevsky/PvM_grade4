package models

type Student struct {
	TelegramId string
	FullName   string
}

func NewStudent(telegramId string, fullName string) Student {
	return Student{
		TelegramId: telegramId,
		FullName:   fullName,
	}
}
