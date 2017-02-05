package slackApi

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"
	"time"
)

var myToken = os.Getenv("VerificationToken")
var templates = template.Must(template.ParseGlob("templates/*.json"))

var netClient = &http.Client {
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
	// user := payload.User
	// server := payload.Actions[0].Value
	// response_url := payload.Response_url
	//
	// _, err := netClient.Post()
	//
	// if err == nil {
	// 	templates.ExecuteTemplate(w, "request_submitted", "")
	// }
	templates.ExecuteTemplate(w, "request_submitted", "")
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
		templates.ExecuteTemplate(w, "coming_soon.json", "")
	default:

		templates.ExecuteTemplate(w, "choose_action.json", "")
	}
}
