package slackApi

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

func NewDB(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func NewRequest(db *sql.DB, requester string) OperationsRequest {
    var id int64

    err := db.QueryRow(
        "insert into Requests (requester) values ($1) returning id", requester,
    ).Scan(&id)

	if err != nil {
		log.Fatal(err)
	}

	return OperationsRequest{
		Id:        id,
		Requester: requester,
	}
}

func LoadRequest(db *sql.DB, requestId string) OperationsRequest {
	var (
		id        int64
		requester string
		action    sql.NullString
		server    sql.NullString
	)

	err := db.QueryRow(
		"select id, requester, action, server from Requests where id = $1", requestId,
	).Scan(&id, &requester, &action, &server)
	if err != nil {
		log.Fatal(err)
	}

	return OperationsRequest{
		Id:        id,
		Requester: requester,
		Action:    action.String,
		Server:    server.String,
	}
}

func UpdateRequest(db *sql.DB, opsRequest OperationsRequest) error {
    query := "update Requests set action = $2, server = $3, responder = $4, approved = $5 where id = $1"
	_, err := db.Exec(query, opsRequest.Id, opsRequest.Action, opsRequest.Server, opsRequest.Responder, opsRequest.Approved)

	if err != nil {
		return err
	}

	return nil
}