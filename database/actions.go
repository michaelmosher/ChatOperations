package database

import (
	"log"

	_ "github.com/go-sql-driver/mysql"

	"chatOperations/operations"
)

type ActionRepo struct {
	*DB
}

func (repo *ActionRepo) FindById(actionId int) (operations.Action, error) {
	var (
		id      int64
		title   string
		command string
	)

	err := repo.QueryRow(
		"select id, title, command from Actions where id = $1", actionId,
	).Scan(&id, &title, &command)

	return operations.Action{
		Id:      id,
		Title:   title,
		Command: command,
	}, err
}

func (repo *ActionRepo) FindAll() ([]operations.Action, error) {
	var results = []operations.Action{}

	rows, err := repo.Query("select id, title, command from Actions")
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

		newAction := operations.Action{Id: id,
			Title:   title,
			Command: command}

		results = append(results, newAction)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return results, err
}
