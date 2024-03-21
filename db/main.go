package db

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"reflect"
	"strings"
)

const dbName = "./data/bot_db.db"

type Controller struct {
	db *sql.DB
}

func New() (Controller, error) {
	if _, err := os.Stat("data"); os.IsNotExist(err) {
		err = os.MkdirAll("data", 0755) // 0755 is the permission mode for the directory
		if err != nil {
			return Controller{}, errors.New(fmt.Sprintf("error creating db directory: %s", err))
		}
	} else if err != nil {
		return Controller{}, errors.New(fmt.Sprintf("error checking db directory: %s", err))
	}

	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return Controller{}, errors.New(fmt.Sprintf("failed to open db connection: %s", err))
	}

	return Controller{
		db: db,
	}, nil
}

func (controller *Controller) MigrateUp() error {
	createDiscipline := `
        CREATE TABLE IF NOT EXISTS disciplines (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL
        );
	`
	_, err := controller.db.Exec(createDiscipline)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to create disciplines table: %s", err))
	}

	createStudent := `
        CREATE TABLE IF NOT EXISTS students (
            telegram_id TEXT PRIMARY KEY,
            full_name TEXT NOT NULL
        );
	`
	_, err = controller.db.Exec(createStudent)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to create students table: %s", err))
	}

	createTeachers := `
        CREATE TABLE IF NOT EXISTS teachers (
            telegram_id TEXT PRIMARY KEY,
            full_name TEXT NOT NULL
        );
	`
	_, err = controller.db.Exec(createTeachers)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to create teachers table: %s", err))
	}

	createDisciplineAttendance := `
        CREATE TABLE IF NOT EXISTS discipline_attendance (
			id INTEGER PRIMARY KEY,
            telegram_id TEXT NOT NULL,
            discipline_id INTEGER NOT NULL
        );
	`
	_, err = controller.db.Exec(createDisciplineAttendance)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to create discipline attendance table: %s", err))
	}

	createTask := `
        CREATE TABLE IF NOT EXISTS tasks (
			id INTEGER PRIMARY KEY,
			discipline_id INTEGER,
            name TEXT NOT NULL,
            description TEXT NOT NULL,
			max_grade INTEGER NOT NULL
        );
	`
	_, err = controller.db.Exec(createTask)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to create tasks table: %s", err))
	}

	createGrade := `
        CREATE TABLE IF NOT EXISTS grades (
			id INTEGER PRIMARY KEY,
            telegram_id TEXT NOT NULL,
			task_id INTEGER NOT NULL,
			grade INTEGER NOT NULL
        );
	`
	_, err = controller.db.Exec(createGrade)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to create grade table: %s", err))
	}

	return nil
}

func (controller *Controller) Insert(tableName string, values map[string]any) error {
	var keysInsert []string
	var valuesInsert []string
	for key, value := range values {
		keysInsert = append(keysInsert, key)
		switch reflect.TypeOf(value).Kind() {
		case reflect.String:
			valuesInsert = append(valuesInsert, fmt.Sprintf("'%s'", value))
		case reflect.Uint64:
			valuesInsert = append(valuesInsert, fmt.Sprintf("%d", value))
		}
	}

	insertSQL := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, strings.Join(keysInsert, ", "),
		strings.Join(valuesInsert, ", "))

	tx, err := controller.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(insertSQL)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (controller *Controller) Select(tableName string, parseFunc func(*sql.Rows) (any, error), whereParams map[string]any, cols ...string) (any, error) {
	var whereParam []string
	for key, value := range whereParams {
		switch reflect.TypeOf(value).Kind() {
		case reflect.String:
			whereParam = append(whereParam, fmt.Sprintf("%s = '%s'", key, value))
		case reflect.Int | reflect.Int32 | reflect.Int64 | reflect.Int16 | reflect.Int8 |
			reflect.Uint | reflect.Uint16 | reflect.Uint32 | reflect.Uint64 | reflect.Uint8:
			whereParam = append(whereParam, fmt.Sprintf("%s = %d", key, value))
		}
	}

	whereString := ""
	if len(whereParam) != 0 {
		whereString = fmt.Sprintf("WHERE %s", strings.Join(whereParam, " AND "))
	}

	rows, err := controller.db.Query(fmt.Sprintf("SELECT %s FROM %s %s", strings.Join(cols, ", "), tableName, whereString))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if err = rows.Err(); err != nil {
		return nil, err
	}

	parsed, err := parseFunc(rows)
	if err != nil {
		return nil, err
	}

	return parsed, nil
}

func (controller *Controller) Update(tableName string, updateParams map[string]any, whereParams map[string]any) error {
	var whereParam []string
	for key, value := range whereParams {
		switch reflect.TypeOf(value).Kind() {
		case reflect.String:
			whereParam = append(whereParam, fmt.Sprintf("%s = '%s'", key, value))
		case reflect.Uint64:
			whereParam = append(whereParam, fmt.Sprintf("%s = %d", key, value))
		}
	}

	whereString := ""
	if len(whereParam) != 0 {
		whereString = fmt.Sprintf("WHERE %s", strings.Join(whereParam, " AND "))
	}

	var updateParam []string
	for key, value := range updateParams {
		switch reflect.TypeOf(value).Kind() {
		case reflect.String:
			updateParam = append(updateParam, fmt.Sprintf("%s = '%s'", key, value))
		case reflect.Uint64:
			updateParam = append(updateParam, fmt.Sprintf("%s = %d", key, value))
		}
	}

	_, err := controller.db.Exec(fmt.Sprintf("UPDATE %s SET %s %s", tableName, strings.Join(updateParam, ", "), whereString))
	if err != nil {
		return err
	}

	return nil
}

func (controller *Controller) Delete(tableName string, whereParams map[string]any) error {
	var whereParam []string
	for key, value := range whereParams {
		switch reflect.TypeOf(value).Kind() {
		case reflect.String:
			whereParam = append(whereParam, fmt.Sprintf("%s = '%s'", key, value))
		case reflect.Uint64:
			whereParam = append(whereParam, fmt.Sprintf("%s = %d", key, value))
		}
	}

	whereString := ""
	if len(whereParam) != 0 {
		whereString = fmt.Sprintf("WHERE %s", strings.Join(whereParam, " AND "))
	}

	_, err := controller.db.Exec(fmt.Sprintf("DELETE FROM %s %s", tableName, whereString))
	if err != nil {
		return err
	}

	return nil
}
