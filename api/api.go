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

var DB, _ = NewDB(os.Getenv("DATABASE_URL"))

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

func reportError(err error, response_url string) {
	// post error to response_url
}

func chooseActionReponse(w http.ResponseWriter, payload SlackPayload, opsRequest OperationsRequest) {
	action := payload.Actions[0].Value

	if action == "else" {
		templates.ExecuteTemplate(w, "coming_soon.json", "")
		return
	}

	templates.ExecuteTemplate(w, "choose_server.json", opsRequest)

	err := SaveAction(DB, opsRequest, action)

	if err != nil {
		response_url := payload.Response_url
		reportError(err, response_url)
	}
}

func chooseServerResponse(w http.ResponseWriter, payload SlackPayload, opsRequest OperationsRequest) {
	server := payload.Actions[0].Value
	opsRequest.Server = server
	// response_url := payload.Response_url

	reader, writer := io.Pipe()

	// writing without a reader will deadlock so write in a goroutine
	go func() {
		// it is important to close the writer or reading from the other end of the
		// pipe will never finish
		defer writer.Close()

		templates.ExecuteTemplate(writer, "ops_request_submitted.json", opsRequest)
	}()

	resp, err := netClient.Post(webhookUrl, "application/json", reader)

	if err != nil {
		templates.ExecuteTemplate(w, "request_not_submitted", err.Error)
		return
	}

	defer resp.Body.Close()

	templates.ExecuteTemplate(w, "request_submitted.json", "")
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
		OpsRequest = NewRequest(DB, payload.User.Name)
		chooseActionReponse(w, payload, OpsRequest)
	case "choose_server":
		OpsRequest = LoadRequest(DB, payload.Callback_id)
		chooseServerResponse(w, payload, OpsRequest)
	case "ops_request_submitted":
		templates.ExecuteTemplate(w, "under_development.json", "")
	default:
		templates.ExecuteTemplate(w, "choose_action.json", "")
	}
}
