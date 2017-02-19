package server

import (
	"encoding/json"
	"html/template"
	"net/http"
)

func New(token string, s *State, n *Notifier) Server {
	templates := template.Must(template.ParseGlob("templates/*.json"))

	return Server{
		verificationToken: token,
		state:             s,
		notifier:          n,
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

func (server *Server) chooseActionResponse(w http.ResponseWriter, p SlackPayload) {
	action := p.Actions[0].Value

	if action == "else" {
		server.jsonTemplates.ExecuteTemplate(w, "coming_soon.json", "")
		return
	}

	id, _ := server.state.SaveRequest(OperationsRequest{
		Request: payload.User.Name,
		Action: action
	})

	server.jsonTemplates.ExecuteTemplate(w, "choose_server.json", id)
}

func (server *Server) chooseServerResponse(w http.ResponseWriter, p SlackPayload) {
	opsRequest := server.state.LoadRequest(p.Callback_id)

	opsRequest.Server = p.Actions[0].Value
	opsRequest.Response_url = p.Response_url

	err := server.notifier.NotifyRequestSubmitted(opsRequest)

	if err != nil {
		server.jsonTemplates.ExecuteTemplate(w, "request_not_submitted.json", err.Error)
		return
	}

	server.jsonTemplates.ExecuteTemplate(w, "request_submitted.json", "")

	err = server.state.UpdateRequest(opsRequest)

	if err != nil {
		server.notifier.reportError(opsRequest.Response_url, err.Error)
	}
}

func (server *Server) opsResponseReceiveResponse(w http.ResponseWriter, p SlackPayload) {
	opsRequest := server.state.LoadRequest(p.Callback_id)
	opsRequest.Responder = p.User.Name

	var err error
	if p.Actions[0].Value == "approved" {
		opsRequest.Approved = true
		server.jsonTemplates.ExecuteTemplate(w, "ops_request_approved.json", opsRequest)
		err = server.notifier.NotifyRequestApproved(opsRequest)
	} else {
		opsRequest.Approved = false
		server.jsonTemplates.ExecuteTemplate(w, "ops_request_rejected.json", opsRequest)
		err = server.notifier.NotifyRequestRejected(opsRequest)
	}

	if err != nil {
		response_url := p.Response_url
		server.notifier.reportError(response_url, err.Error)
	}

	err = server.state.UpdateRequest(opsRequest)

	if err != nil {
		response_url := p.Response_url
		server.notifier.reportError(response_url, err.Error)
	}
	// do thing
}

func (server *Server) routeRequest(w http.ResponseWriter, p SlackPayload) {
	switch p.Actions[0].Name {
	case "choose_action":
		server.chooseActionResponse(w, p)
	case "choose_server":
		server.chooseServerResponse(w, p)
	case "ops_request_submitted":
		server.opsResponseReceiveResponse(w, p)
	default:
		server.jsonTemplates.ExecuteTemplate(w, "choose_action.json", "")
	}
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
