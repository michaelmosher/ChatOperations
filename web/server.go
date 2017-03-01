package web

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"

	"chatoperations/application"
	"chatoperations/operations"
)

var templates = template.Must(template.ParseGlob("web/templates/*.json"))

type Server struct {
	VerificationToken string
	OpsInteractor     application.OperationsInteractor
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
	actionIdString := p.Actions[0].Value

	if actionIdString == "else" {
		return "coming_soon.json", ""
	}

	actionId, _ := strconv.Atoi(actionIdString)
	opsRequest, err := server.OpsInteractor.SetRequestAction(operations.Request{
		Requester:    p.User.Name,
		Response_url: p.Response_url,
	}, actionId)

	if err != nil {
		return "something_went_wrong", err.Error
	}
	return "choose_server.json", opsRequest.Id
}

func (server *Server) chooseServer(p SlackPayload) (string, interface{}) {
	requestId, _ := strconv.Atoi(p.Callback_id)
	serverId, _ := strconv.Atoi(p.Actions[0].Value)

	opsRequest, err := server.OpsInteractor.SetRequestServer(requestId, serverId)

	if err != nil {
		return "something_went_wrong", err.Error
	}
	return "confirm_request.json", opsRequest
}

func (server *Server) submitRequest(p SlackPayload) (string, interface{}) {
	requestId, _ := strconv.Atoi(p.Callback_id)

	if p.Actions[0].Value == "cancel" {
		return "request_not_submitted", ""
	}

	err := server.OpsInteractor.SubmitRequest(requestId)

	if err != nil {
		return "something_went_wrong.json", err.Error
	}

	return "request_submitted.json", ""
}

func (server *Server) opsResponseReceive(p SlackPayload) (string, interface{}) {
	requestId, _ := strconv.Atoi(p.Callback_id)
	responder := p.User.Name

	var (
		templateName string
		opsRequest   operations.Request
	)

	if p.Actions[0].Value == "approved" {
		templateName = "ops_request_approved.json"
		opsRequest, _ = server.OpsInteractor.ApproveRequest(requestId, responder)
	} else {
		templateName = "ops_request_rejected.json"
		opsRequest, _ = server.OpsInteractor.RejectRequest(requestId, responder)
	}

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
	case "submit_request":
		responseTemplate, data = server.submitRequest(p)
	case "ops_request_submitted":
		responseTemplate, data = server.opsResponseReceive(p)
	default:
		responseTemplate = "choose_action.json"
		data = ""
	}

	templates.ExecuteTemplate(w, responseTemplate, data)
}

func (server *Server) Operations(w http.ResponseWriter, r *http.Request) {
	var p SlackPayload
	p.unmarshal(r)

	// make sure actions has at least one thing in it
	p.Actions = append(p.Actions, SlackAction{})

	if requestToken(r, p) != server.VerificationToken {
		http.Error(w, "Invalid Token", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	server.routeRequest(w, p)
}
