package models

type DisciplineAttendance struct {
	Id           uint64
	TelegramId   string
	DisciplineId uint64
}

func NewDisciplineAttendance(id uint64, telegramId string, disciplineId uint64) DisciplineAttendance {
	return DisciplineAttendance{
		Id:           id,
		TelegramId:   telegramId,
		DisciplineId: disciplineId,
	}
}
