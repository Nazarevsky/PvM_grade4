package controllers

import (
	"database/sql"
	"errors"
	"lab4/db"
	"lab4/db/models"
)

const dbDisciplineAttendance string = "discipline_attendance"

type DisciplineAttendance struct {
	dbController db.Controller
}

func NewDisciplineAttendanceController(dbController db.Controller) DisciplineAttendance {
	return DisciplineAttendance{
		dbController,
	}
}

func parseDisciplineAttendance(rows *sql.Rows) (any, error) {
	var values []models.DisciplineAttendance
	for rows.Next() {
		var value models.DisciplineAttendance

		err := rows.Scan(&value.Id, &value.TelegramId, &value.DisciplineId)
		if err != nil {
			return nil, err
		}
		values = append(values, value)
	}

	return values, nil
}

func (dc *DisciplineAttendance) GetById(id uint64) (models.DisciplineAttendance, error) {
	disciplinesAny, err := dc.dbController.Select(dbDisciplineAttendance, parseDisciplineAttendance, map[string]any{"id": id}, "*")
	if err != nil {
		return models.DisciplineAttendance{}, err
	}

	disciplines := disciplinesAny.([]models.DisciplineAttendance)
	if len(disciplines) == 0 {
		return models.DisciplineAttendance{}, errors.New("not found")
	}

	return disciplines[0], nil
}

func (dc *DisciplineAttendance) GetByTgId(tg_id string) ([]models.DisciplineAttendance, error) {
	disciplinesAny, err := dc.dbController.Select(dbDisciplineAttendance, parseDisciplineAttendance, map[string]any{"telegram_id": tg_id}, "*")
	if err != nil {
		return []models.DisciplineAttendance{}, err
	}

	disciplines := disciplinesAny.([]models.DisciplineAttendance)
	if len(disciplines) == 0 {
		return []models.DisciplineAttendance{}, errors.New("not found")
	}

	return disciplines, nil
}

func (dc *DisciplineAttendance) Create(val models.DisciplineAttendance) error {
	return dc.dbController.Insert(dbDisciplineAttendance, map[string]any{"discipline_id": val.DisciplineId,
		"telegram_id": val.TelegramId})
}

func (dc *DisciplineAttendance) UpdateById(id uint64, val models.DisciplineAttendance) error {
	return dc.dbController.Update(dbDisciplineAttendance, map[string]any{"discipline_id": val.DisciplineId, "telegram_id": val.TelegramId},
		map[string]any{"id": id})
}

func (dc *DisciplineAttendance) RemoveById(id uint64) error {
	return dc.dbController.Delete(dbDisciplineAttendance, map[string]any{"id": id})
}
