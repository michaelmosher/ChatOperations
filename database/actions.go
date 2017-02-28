package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"

	"chatoperations/operations"
)

type ActionRepo struct {
	*sql.DB
}

func NewActionRepo(postgresUrl string) (*ActionRepo, error) {
	db, err := sql.Open("postgres", postgresUrl)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &ActionRepo{db}, nil
}

func (repo *ActionRepo) FindById(id string) (operations.Action, error) {
	var (
		id      int64
		title   string
		command string
	)

	err := repo.QueryRow(
		"select id, title, command from Actions where id = $1", id,
	).Scan(&id, &title, &command)

	return operations.Action{
		Id:      id,
		Title:   title,
		Command: command,
	}, err
}

func (repo *ActionRepo) FindAll() ([]operations.Action, error) {
	var results = []operations.Action

	rows, err := db.Query("select id, title, command from Actions")
	if err != nil {
        log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id      int64
			title   string
			command string
		)

		if err := rows.Scan(&id, &title, &command); err != nil {
	        log.Fatal(err)
        }

		results = append(results, operations.Action{
			Id:      id,
			Title:   title,
			Command, command,
		})
	}
	if err := rows.Err(); err != nil {
	        log.Fatal(err)
	}

	return results, err
}
