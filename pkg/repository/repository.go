package repository

import "database/sql"

type Evraz interface {
}

type Repository struct {
	Evraz
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Evraz: NewEvrazPostgres(db),
	}
}
