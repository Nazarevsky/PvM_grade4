package controllers

import (
	"database/sql"
	"errors"
	"lab4/db"
	"lab4/db/models"
)

const dbTask string = "tasks"

type Task struct {
	dbController db.Controller
}

func NewTaskController(dbController db.Controller) Task {
	return Task{
		dbController,
	}
}

func parseTask(rows *sql.Rows) (any, error) {
	var values []models.Task
	for rows.Next() {
		var value models.Task

		err := rows.Scan(&value.Id, &value.Name, &value.Description, &value.MaxGrade)
		if err != nil {
			return nil, err
		}
		values = append(values, value)
	}

	return values, nil
}

func (dc *Task) GetById(id uint64) (models.Task, error) {
	valueAny, err := dc.dbController.Select(dbTask, parseTask, map[string]any{"id": id}, "*")
	if err != nil {
		return models.Task{}, err
	}

	values := valueAny.([]models.Task)
	if len(values) == 0 {
		return models.Task{}, errors.New("not found")
	}

	return values[0], nil
}

func (dc *Task) GetByDisciplineId(discipline_id uint64) ([]models.Task, error) {
	valueAny, err := dc.dbController.Select(dbTask, parseTask, map[string]any{"discipline_id": discipline_id}, "*")
	if err != nil {
		return []models.Task{}, err
	}

	values := valueAny.([]models.Task)
	if len(values) == 0 {
		return []models.Task{}, errors.New("not found")
	}

	return values, nil
}

func (dc *Task) Create(val models.Task) error {
	return dc.dbController.Insert(dbTask, map[string]any{"name": val.Name, "description": val.Description,
		"max_grade": val.MaxGrade})
}

func (dc *Task) UpdateById(id uint64, val models.Task) error {
	return dc.dbController.Update(dbTask, map[string]any{"name": val.Name, "description": val.Description,
		"max_grade": val.MaxGrade},
		map[string]any{"id": id})
}

func (dc *Task) RemoveById(id uint64) error {
	return dc.dbController.Delete(dbTask, map[string]any{"id": id})
}
