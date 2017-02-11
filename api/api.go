package slackApi

import (
	"encoding/json"
	"html/template"
	"io"
	"net/http"
)

var templates = template.Must(template.ParseGlob("templates/*.json"))

func (payload *SlackPayload) unmarshal(r *http.Request) {
	json.Unmarshal([]byte(r.PostFormValue("payload")), &payload)
}

func requestToken(r *http.Request, p SlackPayload) string {
	switch {
	case r.PostFormValue("token") != "":
		// request came from a Slack slash-command
		return r.PostFormValue("token")
	case  p.Token != "":
		// request came from a Slack interactive message
		return p.Token
	default:
		return ""
	}
}

func (env *Env) httpPost(url string, template_name string, opsRequest OperationsRequest) error {
	reader, writer := io.Pipe()

	// writing without a reader will deadlock so write in a goroutine
	go func() {
		defer writer.Close()
		templates.ExecuteTemplate(writer, template_name, opsRequest)
	}()

	resp, err := env.NetClient.Post(url, "application/json", reader)

	if err != nil {
		return err
	}

	defer resp.Body.Close()
	return nil
}

func (env *Env) reportError(err error, response_url string) {
	// post error to response_url
	reader, writer := io.Pipe()

	// writing without a reader will deadlock so write in a goroutine
	go func() {
		defer writer.Close()
		templates.ExecuteTemplate(writer, "something_went_wrong.json", error.Error)
	}()

	resp, _ := env.NetClient.Post(response_url, "application/json", reader)

	defer resp.Body.Close()
}

func (env *Env) chooseActionReponse(w http.ResponseWriter, payload SlackPayload, opsRequest OperationsRequest) {
	opsRequest.Action = payload.Actions[0].Value

	if opsRequest.Action == "else" {
		templates.ExecuteTemplate(w, "coming_soon.json", "")
		return
	}

	templates.ExecuteTemplate(w, "choose_server.json", opsRequest)

	err := env.Db.UpdateRequest(opsRequest)

	if err != nil {
		response_url := payload.Response_url
		env.reportError(err, response_url)
	}
}

func (env *Env) chooseServerResponse(w http.ResponseWriter, payload SlackPayload, opsRequest OperationsRequest) {
	opsRequest.Server = payload.Actions[0].Value
	opsRequest.Response_url = payload.Response_url

	err := env.httpPost(env.WebhookUrl, "ops_request_submitted.json", opsRequest)

	if err != nil {
		templates.ExecuteTemplate(w, "request_not_submitted.json", err.Error)
		return
	}

	templates.ExecuteTemplate(w, "request_submitted.json", "")

	err = env.Db.UpdateRequest(opsRequest)

	if err != nil {
		env.reportError(err, opsRequest.Response_url)
	}
}

func (env *Env) opsResponseReceiveResponse(w http.ResponseWriter, payload SlackPayload, opsRequest OperationsRequest) {
	opsRequest.Responder = payload.User.Name

	var err error
	if payload.Actions[0].Value == "approved" {
		opsRequest.Approved = true
		templates.ExecuteTemplate(w, "ops_request_approved.json", opsRequest)
		err = env.httpPost(opsRequest.Response_url, "request_approved.json", opsRequest)
	} else {
		opsRequest.Approved = false
		templates.ExecuteTemplate(w, "ops_request_rejected.json", opsRequest)
		err = env.httpPost(opsRequest.Response_url, "request_rejected.json", opsRequest)
	}

	if err != nil {
		response_url := payload.Response_url
		env.reportError(err, response_url)
	}

	err = env.Db.UpdateRequest(opsRequest)

	if err != nil {
		response_url := payload.Response_url
		env.reportError(err, response_url)
	}
	// do thing
}

func (env *Env) Operations(w http.ResponseWriter, r *http.Request) {
	var payload SlackPayload
	payload.unmarshal(r)

	if requestToken(r, payload) != env.VerificationToken {
		http.Error(w, "Invalid Token", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if len(payload.Actions) == 0 {
		templates.ExecuteTemplate(w, "choose_action.json", "")
		return
	}

	switch payload.Actions[0].Name {
	case "choose_action":
		OpsRequest := env.Db.NewRequest(payload.User.Name)
		env.chooseActionReponse(w, payload, OpsRequest)
	case "choose_server":
		OpsRequest := env.Db.LoadRequest(payload.Callback_id)
		env.chooseServerResponse(w, payload, OpsRequest)
	case "ops_request_submitted":
		OpsRequest := env.Db.LoadRequest(payload.Callback_id)
		env.opsResponseReceiveResponse(w, payload, OpsRequest)
	default:
		templates.ExecuteTemplate(w, "choose_action.json", "")
	}
}
