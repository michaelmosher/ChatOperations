package server

import (
	"html/template"
	"net/http"
	"time"
)

type State interface {
	SaveRequest(opsRequest OperationsRequest) (id int64, err error)
	LoadRequest(requestId string) (OperationsRequest, error)
}

type Notifier interface {
	NotifyRequestSubmitted(opsRequest OperationsRequest) error
	NotifyRequestApproved(opsRequest OperationsRequest) error
	NotifyRequestRejected(opsRequest OperationsRequest) error
	NotifyError(url string, originalErr error)
}

type Server struct {
	verificationToken string
	state             *State
	notifier          *Notifier
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
