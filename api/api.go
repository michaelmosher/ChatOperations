package slackApi

import (
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"os"
	"time"
)

var myToken = os.Getenv("VerificationToken")
var webhookUrl = os.Getenv("WebhookUrl")
var templates = template.Must(template.ParseGlob("templates/*.json"))

var db, _ = NewDB(os.Getenv("DATABASE_URL"))

var netClient = &http.Client{
	Timeout: time.Second * 10,
}

func authorized(r *http.Request, p SlackPayload) bool {
	// request came from a Slack slash-command
	if r.PostFormValue("token") == myToken {
		return true
	}

	// request came from a Slack interactive message
	if p.Token == myToken {
		return true
	}

	return false
}

func httpPost(url string, template_name string, opsRequest OperationsRequest) error {
	reader, writer := io.Pipe()

	// writing without a reader will deadlock so write in a goroutine
	go func() {
		defer writer.Close()
		templates.ExecuteTemplate(writer, template_name, opsRequest)
	}()

	resp, err := netClient.Post(url, "application/json", reader)

	if err != nil {
		return err
	}

	defer resp.Body.Close()
	return nil
}

func reportError(err error, response_url string) {
	// post error to response_url
	reader, writer := io.Pipe()

	// writing without a reader will deadlock so write in a goroutine
	go func() {
		defer writer.Close()
		templates.ExecuteTemplate(writer, "something_went_wrong.json", error.Error)
	}()

	resp, _ := netClient.Post(response_url, "application/json", reader)

	defer resp.Body.Close()
}

func chooseActionReponse(w http.ResponseWriter, payload SlackPayload, opsRequest OperationsRequest) {
	opsRequest.Action = payload.Actions[0].Value

	if opsRequest.Action == "else" {
		templates.ExecuteTemplate(w, "coming_soon.json", "")
		return
	}

	templates.ExecuteTemplate(w, "choose_server.json", opsRequest)

	err := UpdateRequest(db, opsRequest)

	if err != nil {
		response_url := payload.Response_url
		reportError(err, response_url)
	}
}

func chooseServerResponse(w http.ResponseWriter, payload SlackPayload, opsRequest OperationsRequest) {
	opsRequest.Server = payload.Actions[0].Value
	opsRequest.Response_url = payload.Response_url

	err := httpPost(webhookUrl, "ops_request_submitted.json", opsRequest)

	if err != nil {
		templates.ExecuteTemplate(w, "request_not_submitted.json", err.Error)
		return
	}

	templates.ExecuteTemplate(w, "request_submitted.json", "")

	err = UpdateRequest(db, opsRequest)

	if err != nil {
		reportError(err, opsRequest.Response_url)
	}
}

func opsResponseReceiveResponse(w http.ResponseWriter, payload SlackPayload, opsRequest OperationsRequest) {
	opsRequest.Responder = payload.User.Name

	if payload.Actions[0].Value == "approved" {
		opsRequest.Approved = true
		templates.ExecuteTemplate(w, "ops_request_approved.json", opsRequest)
		err := httpPost(opsRequest.Response_url, "request_approved.json", opsRequest)
	} else {
		opsRequest.Approved = false
		templates.ExecuteTemplate(w, "ops_request_rejected.json", opsRequest)
		err := httpPost(opsRequest.Response_url, "request_rejected.json", opsRequest)
	}

	if err != nil {
		response_url := payload.Response_url
		reportError(err, response_url)
	}

	err = UpdateRequest(db, opsRequest)

	if err != nil {
		response_url := payload.Response_url
		reportError(err, response_url)
	}
	// do thing
}

func Operations(w http.ResponseWriter, r *http.Request) {
	var payload SlackPayload
	json.Unmarshal([]byte(r.PostFormValue("payload")), &payload)

	if !authorized(r, payload) {
		http.Error(w, "Invalid Token", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if len(payload.Actions) == 0 {
		templates.ExecuteTemplate(w, "choose_action.json", "")
		return
	}

	var OpsRequest OperationsRequest

	switch payload.Actions[0].Name {
	case "choose_action":
		OpsRequest = NewRequest(db, payload.User.Name)
		chooseActionReponse(w, payload, OpsRequest)
	case "choose_server":
		OpsRequest = LoadRequest(db, payload.Callback_id)
		chooseServerResponse(w, payload, OpsRequest)
	case "ops_request_submitted":
		OpsRequest = LoadRequest(db, payload.Callback_id)
		opsResponseReceiveResponse(w, payload, OpsRequest)
	default:
		templates.ExecuteTemplate(w, "choose_action.json", "")
	}
}
