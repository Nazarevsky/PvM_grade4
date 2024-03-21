package controllers

import (
	"database/sql"
	"errors"
	"lab4/db"
	"lab4/db/models"
)

const dbDiscipline string = "disciplines"

type Discipline struct {
	dbController db.Controller
}

func NewDisciplineController(dbController db.Controller) Discipline {
	return Discipline{
		dbController,
	}
}

func parseDisciplines(rows *sql.Rows) (any, error) {
	var disciplines []models.Discipline
	for rows.Next() {
		var discipline models.Discipline

		err := rows.Scan(&discipline.Id, &discipline.Name)
		if err != nil {
			return nil, err
		}
		disciplines = append(disciplines, discipline)
	}

	return disciplines, nil
}

func (dc *Discipline) GetById(id uint64) (models.Discipline, error) {
	disciplinesAny, err := dc.dbController.Select(dbDiscipline, parseDisciplines, map[string]any{"id": id}, "*")
	if err != nil {
		return models.Discipline{}, err
	}

	disciplines := disciplinesAny.([]models.Discipline)
	if len(disciplines) == 0 {
		return models.Discipline{}, errors.New("not found")
	}

	return disciplines[0], nil
}

func (dc *Discipline) GetByName(name string) (models.Discipline, error) {
	disciplinesAny, err := dc.dbController.Select(dbDiscipline, parseDisciplines, map[string]any{"name": name}, "*")
	if err != nil {
		return models.Discipline{}, err
	}

	disciplines := disciplinesAny.([]models.Discipline)
	if len(disciplines) == 0 {
		return models.Discipline{}, errors.New("not found")
	}

	return disciplines[0], nil
}

func (dc *Discipline) Create(discipline models.Discipline) error {
	return dc.dbController.Insert(dbDiscipline, map[string]any{"name": discipline.Name})
}

func (dc *Discipline) UpdateById(id uint64, discipline models.Discipline) error {
	return dc.dbController.Update(dbDiscipline, map[string]any{"name": discipline.Name}, map[string]any{"id": id})
}

func (dc *Discipline) RemoveById(id uint64) error {
	return dc.dbController.Delete(dbDiscipline, map[string]any{"id": id})
}
