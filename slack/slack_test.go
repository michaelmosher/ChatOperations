package slack_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"chatoperations/operations"
	"chatoperations/slack"
)

type test_struct struct {
	Text string
}

var slackStub = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Fatal(err)
	}

	var t test_struct
	err = json.Unmarshal(body, &t)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(t.Text)
}))

var slackTest = slack.New(slack.Config{
	TemplatesGlob: "templates/*.json",
	WebhookUrl:    slackStub.URL,
})

var myTestRequest = operations.Request{
	Id:           1,
	Requester:    "tester",
	Action:       "test",
	Server:       "test-server",
	Responder:    "sponder",
	Response_url: slackStub.URL,
}

func TestNotifyRequestSubmitted(t *testing.T) {
	err := slackTest.NotifyRequestSubmitted(myTestRequest)

	if err != nil {
		t.Error("Encountered an error sending a notification: ", err)
	}
}

func TestNotifyRequestApproved(t *testing.T) {
	err := slackTest.NotifyRequestApproved(myTestRequest)

	if err != nil {
		t.Error("Encountered an error sending a notification: ", err)
	}
}

func TestNotifyRequestRejected(t *testing.T) {
	err := slackTest.NotifyRequestRejected(myTestRequest)

	if err != nil {
		t.Error("Encountered an error sending a notification: ", err)
	}
}

func TestNotifyError(t *testing.T) {
	testError := errors.New("Something went wrong")
	slackTest.NotifyError(slackStub.URL, testError)
}
