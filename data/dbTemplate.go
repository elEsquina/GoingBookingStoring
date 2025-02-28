package data

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

type DBTemplate struct {
	db *sqlx.DB
}

func NewDBTemplate(connStr string) *DBTemplate {
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	return &DBTemplate{db: db}
}

func QueryStructs[T any](template *DBTemplate, query string, args ...any) ([]T, error) {
	var results []T
	if err := template.db.Select(&results, query, args...); err != nil {
		return nil, err
	}
	return results, nil
}

func QueryStruct[T any](template *DBTemplate, query string, args ...any) (*T, error) {
	var result T
	if err := template.db.Get(&result, query, args...); err != nil {
		return nil, err
	}
	return &result, nil
}

func ExecuteInsert(template *DBTemplate, query string, args ...any) (int, error) {
	var id int
	err := template.db.QueryRow(query, args...).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}


func ExecuteUpdateOrDelete(template *DBTemplate, query string, args ...any) (int, error) {
	result, err := template.db.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	affectedRows, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(affectedRows), nil
}
