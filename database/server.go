package database

import (
	"log"

	_ "github.com/go-sql-driver/mysql"

	"chatOperations/operations"
)

type ServerRepo struct {
	*DB
}

func (repo *ServerRepo) FindById(serviceId int) (operations.Server, error) {
	var (
		id          int64
		title       string
		address     string
		environment string
	)

	err := repo.QueryRow(
		"select id, title, address, environment from Servers where id = ?", serviceId,
	).Scan(&id, &title, &address, &environment)

	return operations.Server{
		Id:          id,
		Title:       title,
		Address:     address,
		Environment: environment,
	}, err
}

func (repo *ServerRepo) FindAll() ([]operations.Server, error) {
	var results = []operations.Server{}

	rows, err := repo.Query("select id, title, address from Servers")
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

		newServer := operations.Server{Id: id,
			Title:   title,
			Address: address}

		results = append(results, newServer)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return results, err
}
