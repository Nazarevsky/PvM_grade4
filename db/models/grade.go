package models

type Grade struct {
	Id         uint64
	TelegramId string
	TaskId     uint64
	Grade      uint64
}

func NewGrade(id uint64, telegramId string, taskId uint64, grade uint64) Grade {
	return Grade{
		Id:         id,
		TelegramId: telegramId,
		TaskId:     taskId,
		Grade:      grade,
	}
}
