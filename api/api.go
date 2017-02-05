package slackApi

import (
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"os"
	"time"
	"fmt"
	"io/ioutil"
)

var myToken = os.Getenv("VerificationToken")
var webhookUrl = os.Getenv("WebhookUrl")
var templates = template.Must(template.ParseGlob("templates/*.json"))

var netClient = &http.Client{
	Timeout: time.Second * 10,
}

type SlackAction struct {
	Name  string
	Value string
}

type SlackUser struct {
	Id   string
	Name string
}

type SlackChannel struct {
	Id   string
	Name string
}

type SlackTeam struct {
	Id     string
	Domain string
}

type SlackPayload struct {
	Token        string
	Actions      []SlackAction
	Team         SlackTeam
	Channel      SlackChannel
	User         SlackUser
	Callback_id  string
	Message_ts   string
	Response_url string
}

type OpsRequest struct {
	User   SlackUser
	Server string
	Action string
}

func authorized(r *http.Request, p SlackPayload) bool {
	if r.PostFormValue("token") == myToken {
		return true
	}

	if p.Token == myToken {
		return true
	}

	return false
}

func chooseActionReponse(w http.ResponseWriter, payload SlackPayload) {
	if payload.Actions[0].Value == "configLoad" {
		templates.ExecuteTemplate(w, "choose_server.json", "")
	} else {
		templates.ExecuteTemplate(w, "coming_soon.json", "")
	}
}

func chooseServerResponse(w http.ResponseWriter, payload SlackPayload) {
	opsRequest := OpsRequest{
		payload.User,
		payload.Actions[0].Value,
		"config load",
	}
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

	bod, _ := ioutil.ReadAll(resp.Body)
	fmt.Fprintln(w, string(bod))
	// templates.ExecuteTemplate(w, "request_submitted.json", "")
}

func Operations(w http.ResponseWriter, r *http.Request) {
	var payload SlackPayload
	json.Unmarshal([]byte(r.PostFormValue("payload")), &payload)

	if !authorized(r, payload) {
		http.Error(w, "Invalid Token", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	switch payload.Callback_id {
	case "choose_action":
		chooseActionReponse(w, payload)
	case "choose_server":
		chooseServerResponse(w, payload)
	case "ops_request_submitted":
		templates.ExecuteTemplate(w, "under_development.json", "")
	default:

		templates.ExecuteTemplate(w, "choose_action.json", "")
	}
}
