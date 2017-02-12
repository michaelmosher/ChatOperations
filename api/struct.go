package slackApi

import (
	"html/template"
	"net/http"
	"time"
)

type ApiConfig struct {
	VerificationToken string
	WebhookUrl        string
	DatabaseUrl       string
	NetClientTimeout  time.Duration
}

type Api struct {
	VerificationToken string
	WebhookUrl        string
	DB                *DB
	NetClient         *http.Client
	JsonTemplates     *template.Template
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

type OperationsRequest struct {
	Id           int64
	Requester    string
	Server       string
	Action       string
	Responder    string
	Approved     bool
	Response_url string
}
