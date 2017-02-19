package database

import (
	"database/sql"

	_ "github.com/lib/pq"

	"chatoperations/operations"
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

func (db *DB) newRequest(o operations.Request) (id int64, err error) {
	err = db.QueryRow(
		"insert into Requests (requester, action) values ($1, $2) returning id",
		o.Requester, o.Action,
	).Scan(&id)

	return id, err
}

func (db *DB) updateRequest(o operations.Request) (id int64, err error) {
	err = db.QueryRow(
		"update Requests set action = $2, server = $3, responder = $4, approved = $5, response_url = $6 where id = $1 returning id",
		o.Id, o.Action, o.Server, o.Responder, o.Approved, o.Response_url,
	).Scan(&id)

	return id, err
}

func (db *DB) SaveRequest(o operations.Request) (id int64, err error) {
	if o.Id == 0 {
		return db.newRequest(o)
	} else {
		return db.updateRequest(o)
	}
}

func (db *DB) LoadRequest(requestId string) (operations.Request, error) {
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

	return operations.Request{
		Id:           id,
		Requester:    requester,
		Action:       action.String,
		Server:       server.String,
		Responder:    responder.String,
		Approved:     approved.Bool,
		Response_url: response_url.String,
	}, err
}
