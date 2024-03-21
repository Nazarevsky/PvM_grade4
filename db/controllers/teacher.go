package controllers

import (
	"database/sql"
	"errors"
	"lab4/db"
	"lab4/db/models"
)

const dbTeacher string = "teachers"

type Teacher struct {
	dbController db.Controller
}

func NewTeacherController(dbController db.Controller) Teacher {
	return Teacher{
		dbController,
	}
}

func parseTeacher(rows *sql.Rows) (any, error) {
	var values []models.Teacher
	for rows.Next() {
		var value models.Teacher

		err := rows.Scan(&value.TelegramId, &value.FullName)
		if err != nil {
			return nil, err
		}
		values = append(values, value)
	}

	return values, nil
}

func (dc *Teacher) GetById(telegramId string) (models.Teacher, error) {
	valueAny, err := dc.dbController.Select(dbTeacher, parseTeacher, map[string]any{"telegram_id": telegramId}, "*")
	if err != nil {
		return models.Teacher{}, err
	}

	values := valueAny.([]models.Teacher)
	if len(values) == 0 {
		return models.Teacher{}, errors.New("not found")
	}

	return values[0], nil
}

func (dc *Teacher) GetByFullName(fullName string) (models.Teacher, error) {
	valueAny, err := dc.dbController.Select(dbTeacher, parseTeacher, map[string]any{"full_name": fullName}, "*")
	if err != nil {
		return models.Teacher{}, err
	}

	values := valueAny.([]models.Teacher)
	if len(values) == 0 {
		return models.Teacher{}, errors.New("not found")
	}

	return values[0], nil
}

func (dc *Teacher) Create(val models.Teacher) error {
	return dc.dbController.Insert(dbTeacher, map[string]any{"telegram_id": val.TelegramId, "full_name": val.FullName})
}

func (dc *Teacher) UpdateById(telegramId string, val models.Teacher) error {
	return dc.dbController.Update(dbTeacher, map[string]any{"full_name": val.FullName},
		map[string]any{"telegram_id": telegramId})
}

func (dc *Teacher) RemoveById(telegramId string) error {
	return dc.dbController.Delete(dbTeacher, map[string]any{"telegram_id": telegramId})
}
