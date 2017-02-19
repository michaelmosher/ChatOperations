package chatoperations

import (
	"log"

	"database/sql"
	_ "github.com/lib/pq"
)

type Datastore interface {
	NewRequest(requester string) OperationsRequest
	LoadRequest(requestId string) OperationsRequest
	UpdateRequest(opsRequest OperationsRequest) error
}

type DB struct {
	*sql.DB
}

func NewDB(dataSourceName string) (*DB, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

func (db *DB) NewRequest(requester string, action string) OperationsRequest {
	var id int64

	err := db.QueryRow(
		"insert into Requests (requester, action) values ($1, $2) returning id",
		requester, action,
	).Scan(&id)

	if err != nil {
		log.Fatal(err)
	}

	return OperationsRequest{
		Id:        id,
		Requester: requester,
		Action:    action,
	}
}

func (db *DB) LoadRequest(requestId string) OperationsRequest {
	var (
		id           int64
		requester    string
		action       sql.NullString
		server       sql.NullString
		responder    sql.NullString
		approved     sql.NullBool
		response_url sql.NullString
	)

	err := db.QueryRow(
		"select id, requester, action, server, responder, approved, response_url from Requests where id = $1",
		requestId,
	).Scan(&id, &requester, &action, &server, &responder, &approved, &response_url)
	if err != nil {
		log.Fatal(err)
	}

	return OperationsRequest{
		Id:           id,
		Requester:    requester,
		Action:       action.String,
		Server:       server.String,
		Responder:    responder.String,
		Approved:     approved.Bool,
		Response_url: response_url.String,
	}
}

func (db *DB) UpdateRequest(opsRequest OperationsRequest) error {
	query := "update Requests set action = $2, server = $3, responder = $4, approved = $5, response_url = $6 where id = $1"
	_, err := db.Exec(query, opsRequest.Id, opsRequest.Action, opsRequest.Server, opsRequest.Responder, opsRequest.Approved, opsRequest.Response_url)

	if err != nil {
		return err
	}

	return nil
}
