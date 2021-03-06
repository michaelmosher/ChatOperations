package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"chatOperations/operations"
)

type RequestRepo struct {
	*DB
}

func nullIntHelper(id int64) sql.NullInt64 {
	nullInt := sql.NullInt64{Int64: id, Valid: false}

	if id != 0 {
		nullInt.Valid = true
	}

	return nullInt
}

func (repo *RequestRepo) new() (id int64, err error) {
	stmt, err := repo.Prepare("INSERT into Requests SET requester='pending'")
	res, err := stmt.Exec()

	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func (repo *RequestRepo) update(o operations.Request) (id int64, err error) {
	actionId := nullIntHelper(o.Action.Id)
	serverId := nullIntHelper(o.Server.Id)

	stmt, _ := repo.Prepare("update Requests set requester = ?, actionId = ?, serverId = ?, responder = ?, approved = ?, response_url = ? where id = ?")
	res, err := stmt.Exec(o.Requester, actionId, serverId, o.Responder, o.Approved, o.Response_url, o.Id)

	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func (repo *RequestRepo) Store(o operations.Request) (int64, error) {
	if o.Id == 0 {
		return repo.new()
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
		"select id, requester, actionId, serverId, responder, approved, response_url from Requests where id = ?",
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
