package models

type Teacher struct {
	TelegramId string
	FullName   string
}

func NewTeacher(telegramId string, fullName string) Teacher {
	return Teacher{
		TelegramId: telegramId,
		FullName:   fullName,
	}
}
