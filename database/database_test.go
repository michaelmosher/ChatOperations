package database_test

import (
	"strconv"
	"testing"

	"chatoperations/database"
	"chatoperations/operations"
)

var myTestDB, _ = database.New("postgres://michael:@localhost:5432/chatoperations_test")
var myTestRequestId int64

// test db.SaveRequest (new)
func TestNewRequest(t *testing.T) {
	id, err := myTestDB.SaveRequest(operations.Request{
		Requester: "tester",
		Action:    "test",
	})

	if err != nil {
		t.Error("Encountered an error creating a new request: ", err)
	}

	if id == 0 {
		t.Error("Encountered an error creating a new request.")
	}
	myTestRequestId = id
}

func TestLoadRequest(t *testing.T) {
	myTestRequest := operations.Request{
		Id:        myTestRequestId,
		Requester: "tester",
		Action:    "test",
	}

	idString := strconv.FormatInt(myTestRequestId, 10)
	opsRequest, err := myTestDB.LoadRequest(idString)

	if err != nil {
		t.Error("Encountered an error loading a request: ", err)
	}

	if myTestRequest != opsRequest {
		t.Error("Encountered ", opsRequest, " to equal ", myTestRequest)
	}
}

// test db.SaveRequest (update)
func TestUpdateRequest(t *testing.T) {
	myTestRequest := operations.Request{
		Id:        myTestRequestId,
		Requester: "tester",
		Action:    "test",
		Server:    "test-server",
	}

	_, err := myTestDB.SaveRequest(myTestRequest)

	if err != nil {
		t.Error("Encountered an error updating a request: ", err)
	}

	idString := strconv.FormatInt(myTestRequestId, 10)
	opsRequest, err := myTestDB.LoadRequest(idString)

	if err != nil {
		t.Error("Encountered an error loading a request: ", err)
	}

	if myTestRequest != opsRequest {
		t.Error("Expected ", opsRequest, " to equal ", myTestRequest)
	}
}
