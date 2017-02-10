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
		id           int64
		requester    string
		action       sql.NullString
		server       sql.NullString
    	responder    sql.NullString
    	approved     bool
        response_url sql.NullString
	)

	err := db.QueryRow(
		"select id, requester, action, server, request_url from Requests where id = $1", requestId,
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
        Approved:     approved,
        Response_url: response_url.String,
	}
}

func UpdateRequest(db *sql.DB, opsRequest OperationsRequest) error {
    query := "update Requests set action = $2, server = $3, responder = $4, approved = $5, response_url = $6 where id = $1"
	_, err := db.Exec(query, opsRequest.Id, opsRequest.Action, opsRequest.Server, opsRequest.Responder, opsRequest.Approved, opsRequest.Response_url)

	if err != nil {
		return err
	}

	return nil
}
