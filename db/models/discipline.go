package models

type Discipline struct {
	Id   uint64
	Name string
}

func NewDiscipline(id uint64, name string) Discipline {
	return Discipline{
		Id:   id,
		Name: name,
	}
}
