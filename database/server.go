package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"

	"chatoperations/operations"
)

type ServerRepo struct {
	*sql.DB
}

func NewServerepo(postgresUrl string) (*ServerRepo, error) {
	db, err := sql.Open("postgres", postgresUrl)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &ServerRepo{db}, nil
}

func (repo *ServerRepo) FindById(id string) (operations.Server, error) {
	var (
		id          int64
		title       string
		address     string
		environment string
	)

	err := repo.QueryRow(
		"select id, title, address, environment from Servers where id = $1", id,
	).Scan(&id, &title, &address, &environment)

	return operations.Server{
		Id:          id,
		Title:       title,
		Address:     address,
		Environment: environment,
	}, err
}

func (repo *ServerRepo) FindAll() ([]operations.Server, error) {
	var results = []operations.Server

	rows, err := db.Query("select id, title, address from Servers")
	if err != nil {
        log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id      int64
			title   string
			address string
		)

		if err := rows.Scan(&id, &title, &address); err != nil {
	        log.Fatal(err)
        }

		results = append(results, operations.Server{
			Id:      id,
			Title:   title,
			Address, address,
		})
	}
	if err := rows.Err(); err != nil {
	        log.Fatal(err)
	}

	return results, err
}
