package database

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

func New(dataSourceName string) (*DB, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

func (db *DB) NewActionRepo() *ActionRepo {
	return &ActionRepo{db}
}

func (db *DB) NewServerRepo() *ServerRepo {
	return &ServerRepo{db}
}

func (db *DB) NewRequestRepo() *RequestRepo {
	return &RequestRepo{db}
}
