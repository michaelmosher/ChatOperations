package server

import (
	"encoding/json"
	"html/template"
	"net/http"

	"chatoperations/operations"
)

func New(cfg Config) Server {
	templates := template.Must(template.ParseGlob(cfg.TemplatesGlob))

	return Server{
		verificationToken: cfg.VerificationToken,
		state:             cfg.State,
		notifier:          cfg.Notifier,
		jsonTemplates:     templates,
	}
}

func (p *SlackPayload) unmarshal(r *http.Request) {
	json.Unmarshal([]byte(r.PostFormValue("payload")), &p)
}

func requestToken(r *http.Request, p SlackPayload) string {
	if r.PostFormValue("token") != "" {
		// request came from a Slack slash-command
		return r.PostFormValue("token")
	}
	// request came from a Slack interactive message
	return p.Token
}

func (server *Server) chooseAction(p SlackPayload) (string, interface{}) {
	action := p.Actions[0].Value

	if action == "else" {
		return "coming_soon.json", ""
	}

	id, err := server.state.SaveRequest(operations.Request{
		Requester: p.User.Name,
		Action:    action,
	})

	if err != nil {
		return "something_went_wrong", err.Error
	}
	return "choose_server.json", id
}

func (server *Server) chooseServer(p SlackPayload) (string, interface{}) {
	opsRequest, err := server.state.LoadRequest(p.Callback_id)

	if err != nil {
		return "something_went_wrong", err.Error
	}

	opsRequest.Server = p.Actions[0].Value
	opsRequest.Response_url = p.Response_url

	err = server.notifier.NotifyRequestSubmitted(opsRequest)

	if err != nil {
		return "request_not_submitted.json", err.Error
	}

	go func() {
		_, err = server.state.SaveRequest(opsRequest)
		if err != nil {
			server.notifier.NotifyError(opsRequest.Response_url, err)
		}
	}()

	return "request_submitted.json", ""
}

func (server *Server) opsResponseReceive(p SlackPayload) (string, interface{}) {
	opsRequest, err := server.state.LoadRequest(p.Callback_id)

	if err != nil {
		return "something_went_wrong", err.Error
	}

	opsRequest.Responder = p.User.Name

	var templateName string
	if p.Actions[0].Value == "approved" {
		opsRequest.Approved = true
		templateName = "ops_request_approved.json"
		err = server.notifier.NotifyRequestApproved(opsRequest)
	} else {
		opsRequest.Approved = false
		templateName = "ops_request_rejected.json"
		err = server.notifier.NotifyRequestRejected(opsRequest)
	}

	if err != nil {
		response_url := p.Response_url
		server.notifier.NotifyError(response_url, err)
	}

	go func() {
		_, err = server.state.SaveRequest(opsRequest)
		if err != nil {
			server.notifier.NotifyError(p.Response_url, err)
		}
	}()

	// do thing

	return templateName, opsRequest
}

func (server *Server) routeRequest(w http.ResponseWriter, p SlackPayload) {
	var responseTemplate string
	var data interface{}

	switch p.Actions[0].Name {
	case "choose_action":
		responseTemplate, data = server.chooseAction(p)
	case "choose_server":
		responseTemplate, data = server.chooseServer(p)
	case "ops_request_submitted":
		responseTemplate, data = server.opsResponseReceive(p)
	default:
		responseTemplate = "choose_action.json"
		data = ""
	}

	server.jsonTemplates.ExecuteTemplate(w, responseTemplate, data)
}

func (server *Server) Operations(w http.ResponseWriter, r *http.Request) {
	var p SlackPayload
	p.unmarshal(r)

	// make sure actions has at least one thing in it
	p.Actions = append(p.Actions, SlackAction{})

	if requestToken(r, p) != server.verificationToken {
		http.Error(w, "Invalid Token", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	server.routeRequest(w, p)
}
