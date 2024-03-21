package controllers

import (
	"database/sql"
	"errors"
	"lab4/db"
	"lab4/db/models"
)

const dbGrade string = "grades"

type Grade struct {
	dbController db.Controller
}

func NewGradeController(dbController db.Controller) Grade {
	return Grade{
		dbController,
	}
}

func parseGrade(rows *sql.Rows) (any, error) {
	var values []models.Grade
	for rows.Next() {
		var value models.Grade

		err := rows.Scan(&value.Id, &value.TelegramId, &value.TaskId, &value.Grade)
		if err != nil {
			return nil, err
		}
		values = append(values, value)
	}

	return values, nil
}

func (dc *Grade) GetById(id uint64) (models.Grade, error) {
	valueAny, err := dc.dbController.Select(dbGrade, parseGrade, map[string]any{"id": id}, "*")
	if err != nil {
		return models.Grade{}, err
	}

	values := valueAny.([]models.Grade)
	if len(values) == 0 {
		return models.Grade{}, errors.New("not found")
	}

	return values[0], nil
}

func (dc *Grade) Create(val models.Grade) error {
	return dc.dbController.Insert(dbGrade, map[string]any{"task_id": val.TaskId,
		"telegram_id": val.TelegramId, "grade": val.Grade})
}

func (dc *Grade) UpdateById(id uint64, val models.Grade) error {
	return dc.dbController.Update(dbGrade, map[string]any{"task_id": val.TaskId,
		"telegram_id": val.TelegramId, "grade": val.Grade},
		map[string]any{"id": id})
}

func (dc *Grade) RemoveById(id uint64) error {
	return dc.dbController.Delete(dbGrade, map[string]any{"id": id})
}
