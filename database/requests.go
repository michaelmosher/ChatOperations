package database

import (
	"database/sql"

	_ "github.com/lib/pq"

	"chatoperations/operations"
)

type RequestRepo struct {
	*DB
}

func (repo *RequestRepo) new(o operations.Request) (id int64, err error) {
	err = repo.QueryRow(
		"insert into Requests (requester, actionId) values ($1, $2) returning id",
		o.Requester, o.Action.Id,
	).Scan(&id)

	return id, err
}

func (repo *RequestRepo) update(o operations.Request) (id int64, err error) {
	err = repo.QueryRow(
		"update Requests set actionId = $2, serverId = $3, responder = $4, approved = $5, response_url = $6 where id = $1 returning id",
		o.Id, o.Action.Id, o.Server.Id, o.Responder, o.Approved, o.Response_url,
	).Scan(&id)

	return id, err
}

func (repo *RequestRepo) Store(o operations.Request) (int64, error) {
	if o.Id == 0 {
		return repo.new(o)
	}

	return repo.update(o)
}

func (repo *RequestRepo) FindById(requestId int) (operations.Request, error) {
	var (
		id           int64
		requester    string
		actionId     sql.NullInt64
		serverId     sql.NullInt64
		responder    sql.NullString
		approved     sql.NullBool
		response_url sql.NullString
	)

	err := repo.QueryRow(
		"select id, requester, actionId, serverId, responder, approved, response_url from Requests where id = $1",
		requestId,
	).Scan(&id, &requester, &actionId, &serverId, &responder, &approved, &response_url)

	return operations.Request{
		Id:           id,
		Requester:    requester,
		Action:       operations.Action{Id: actionId.Int64},
		Server:       operations.Server{Id: serverId.Int64},
		Responder:    responder.String,
		Approved:     approved.Bool,
		Response_url: response_url.String,
	}, err
}
