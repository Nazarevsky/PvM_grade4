package models

type Task struct {
	Id           uint64
	DisciplineId uint64
	Name         string
	Description  string
	MaxGrade     uint64
}

func NewTask(id uint64, disciplineId uint64, name string, description string, maxGrade uint64) Task {
	return Task{
		Id:           id,
		DisciplineId: disciplineId,
		Name:         name,
		Description:  description,
		MaxGrade:     maxGrade,
	}
}
