package controllers

import (
	"database/sql"
	"errors"
	"lab4/db"
	"lab4/db/models"
)

const dbStudent string = "students"

type Student struct {
	dbController db.Controller
}

func NewStudentController(dbController db.Controller) Student {
	return Student{
		dbController,
	}
}

func parseStudent(rows *sql.Rows) (any, error) {
	var values []models.Student
	for rows.Next() {
		var value models.Student

		err := rows.Scan(&value.TelegramId, &value.FullName)
		if err != nil {
			return nil, err
		}
		values = append(values, value)
	}

	return values, nil
}

func (dc *Student) GetById(telegramId string) (models.Student, error) {
	valueAny, err := dc.dbController.Select(dbStudent, parseStudent, map[string]any{"telegram_id": telegramId}, "*")
	if err != nil {
		return models.Student{}, err
	}

	values := valueAny.([]models.Student)
	if len(values) == 0 {
		return models.Student{}, errors.New("not found")
	}

	return values[0], nil
}

func (dc *Student) GetByFullName(fullName string) (models.Student, error) {
	valueAny, err := dc.dbController.Select(dbStudent, parseStudent, map[string]any{"full_name": fullName}, "*")
	if err != nil {
		return models.Student{}, err
	}

	values := valueAny.([]models.Student)
	if len(values) == 0 {
		return models.Student{}, errors.New("not found")
	}

	return values[0], nil
}

func (dc *Student) Create(val models.Student) error {
	return dc.dbController.Insert(dbStudent, map[string]any{"telegram_id": val.TelegramId, "full_name": val.FullName})
}

func (dc *Student) UpdateById(telegramId string, val models.Student) error {
	return dc.dbController.Update(dbStudent, map[string]any{"full_name": val.FullName},
		map[string]any{"telegram_id": telegramId})
}

func (dc *Student) RemoveById(telegramId string) error {
	return dc.dbController.Delete(dbStudent, map[string]any{"telegram_id": telegramId})
}
