package database

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Datastore interface {
	newRequest(opReq OperationsRequest) (id int64, err error)
	updateRequest(opsRequest OperationsRequest) (id int64, err error)
	SaveRequest(opsRequest OperationsRequest) (id int64, err error)
	LoadRequest(requestId string) (OperationsRequest, error)
}

type DB struct {
	*sql.DB
}

type OperationsRequest struct {
	Id           int64
	Requester    string
	Server       string
	Action       string
	Responder    string
	Approved     bool
	Response_url string
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

func (db *DB) newRequest(opsRequest OperationsRequest) (id int64, err error) {
	err = db.QueryRow(
		"insert into Requests (requester, action) values ($1, $2) returning id",
		opsRequest.Requester, opsRequest.Action,
	).Scan(&id)

	return id, err
}

func (db *DB) updateRequest(opsRequest OperationsRequest) (id int64, err error) {
	err = db.QueryRow(
		"update Requests set action = $2, server = $3, responder = $4, approved = $5, response_url = $6 where id = $1 returning id",
		opsRequest.Id, opsRequest.Action, opsRequest.Server, opsRequest.Responder, opsRequest.Approved, opsRequest.Response_url,
	).Scan(&id)

	return id, err
}

func (db *DB) SaveRequest(opsRequest OperationsRequest) (id int64, err error) {
	if opsRequest.Id == 0 {
		return db.newRequest(opsRequest)
	} else {
		return db.updateRequest(opsRequest)
	}
}

func (db *DB) LoadRequest(requestId string) (OperationsRequest, error) {
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

	return OperationsRequest{
		Id:           id,
		Requester:    requester,
		Action:       action.String,
		Server:       server.String,
		Responder:    responder.String,
		Approved:     approved.Bool,
		Response_url: response_url.String,
	}, err
}
