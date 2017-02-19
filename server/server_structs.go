package server

import (
	"html/template"

	"chatoperations/operations"
)

type State interface {
	SaveRequest(o operations.Request) (id int64, err error)
	LoadRequest(requestId string) (operations.Request, error)
}

type Notifier interface {
	NotifyRequestSubmitted(o operations.Request) error
	NotifyRequestApproved(o operations.Request) error
	NotifyRequestRejected(o operations.Request) error
	NotifyError(url string, originalErr error)
}

type Config struct {
	VerificationToken string
	State             State
	Notifier          Notifier
	TemplatesGlob     string
}

type Server struct {
	verificationToken string
	state             State
	notifier          Notifier
	jsonTemplates     *template.Template
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
