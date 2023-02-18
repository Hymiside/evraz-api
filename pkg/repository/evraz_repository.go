package repository

import "database/sql"

type EvrazPostgres struct {
	db *sql.DB
}

func NewEvrazPostgres(db *sql.DB) *EvrazPostgres {
	return &EvrazPostgres{db: db}
}
