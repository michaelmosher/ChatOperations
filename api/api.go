package slackApi

import (
	"encoding/json"
	"html/template"
	"io"
	"net/http"
)

func New(cfg ApiConfig) Api {
	// connect db
	db, _ := NewDB(cfg.DatabaseUrl)
	// create netClient
	netClient := &http.Client{
		Timeout: cfg.NetClientTimeout,
	}
	// load templates
	templates := template.Must(template.ParseGlob("templates/*.json"))
	// create message Channel
	// start reading

	return Api{
		VerificationToken: cfg.VerificationToken,
		WebhookUrl:        cfg.WebhookUrl,
		DB:                db,
		NetClient:         netClient,
		JsonTemplates:     templates,
	}
}

func (payload *SlackPayload) unmarshal(r *http.Request) {
	json.Unmarshal([]byte(r.PostFormValue("payload")), &payload)
}

func requestToken(r *http.Request, p SlackPayload) string {
	switch {
	case r.PostFormValue("token") != "":
		// request came from a Slack slash-command
		return r.PostFormValue("token")
	case p.Token != "":
		// request came from a Slack interactive message
		return p.Token
	default:
		return ""
	}
}

func (api *Api) httpPost(url string, template_name string, data interface{}) error {
	reader, writer := io.Pipe()

	// writing without a reader will deadlock so write in a goroutine
	go func() {
		defer writer.Close()
		api.JsonTemplates.ExecuteTemplate(writer, template_name, data)
	}()

	resp, err := api.NetClient.Post(url, "application/json", reader)

	if err != nil {
		return err
	}

	defer resp.Body.Close()
	return nil
}

func (api *Api) reportError(err error, response_url string) {
	// post error to response_url
	reader, writer := io.Pipe()

	// writing without a reader will deadlock so write in a goroutine
	go func() {
		defer writer.Close()
		api.JsonTemplates.ExecuteTemplate(writer, "something_went_wrong.json", error.Error)
	}()

	resp, _ := api.NetClient.Post(response_url, "application/json", reader)

	defer resp.Body.Close()
}

func (api *Api) chooseActionReponse(w http.ResponseWriter, payload SlackPayload, opsRequest OperationsRequest) {
	opsRequest.Action = payload.Actions[0].Value

	if opsRequest.Action == "else" {
		api.JsonTemplates.ExecuteTemplate(w, "coming_soon.json", "")
		return
	}

	api.JsonTemplates.ExecuteTemplate(w, "choose_server.json", opsRequest)

	err := api.DB.UpdateRequest(opsRequest)

	if err != nil {
		response_url := payload.Response_url
		api.reportError(err, response_url)
	}
}

func (api *Api) chooseServerResponse(w http.ResponseWriter, payload SlackPayload, opsRequest OperationsRequest) {
	opsRequest.Server = payload.Actions[0].Value
	opsRequest.Response_url = payload.Response_url

	err := api.httpPost(api.WebhookUrl, "ops_request_submitted.json", opsRequest)

	if err != nil {
		api.JsonTemplates.ExecuteTemplate(w, "request_not_submitted.json", err.Error)
		return
	}

	api.JsonTemplates.ExecuteTemplate(w, "request_submitted.json", "")

	err = api.DB.UpdateRequest(opsRequest)

	if err != nil {
		api.reportError(err, opsRequest.Response_url)
	}
}

func (api *Api) opsResponseReceiveResponse(w http.ResponseWriter, payload SlackPayload, opsRequest OperationsRequest) {
	opsRequest.Responder = payload.User.Name

	var err error
	if payload.Actions[0].Value == "approved" {
		opsRequest.Approved = true
		api.JsonTemplates.ExecuteTemplate(w, "ops_request_approved.json", opsRequest)
		err = api.httpPost(opsRequest.Response_url, "request_approved.json", opsRequest)
	} else {
		opsRequest.Approved = false
		api.JsonTemplates.ExecuteTemplate(w, "ops_request_rejected.json", opsRequest)
		err = api.httpPost(opsRequest.Response_url, "request_rejected.json", opsRequest)
	}

	if err != nil {
		response_url := payload.Response_url
		api.reportError(err, response_url)
	}

	err = api.DB.UpdateRequest(opsRequest)

	if err != nil {
		response_url := payload.Response_url
		api.reportError(err, response_url)
	}
	// do thing
}

func (api *Api) Operations(w http.ResponseWriter, r *http.Request) {
	var payload SlackPayload
	payload.unmarshal(r)

	if requestToken(r, payload) != api.VerificationToken {
		http.Error(w, "Invalid Token", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if len(payload.Actions) == 0 {
		api.JsonTemplates.ExecuteTemplate(w, "choose_action.json", "")
		return
	}

	switch payload.Actions[0].Name {
	case "choose_action":
		OpsRequest := api.DB.NewRequest(payload.User.Name)
		api.chooseActionReponse(w, payload, OpsRequest)
	case "choose_server":
		OpsRequest := api.DB.LoadRequest(payload.Callback_id)
		api.chooseServerResponse(w, payload, OpsRequest)
	case "ops_request_submitted":
		OpsRequest := api.DB.LoadRequest(payload.Callback_id)
		api.opsResponseReceiveResponse(w, payload, OpsRequest)
	default:
		api.JsonTemplates.ExecuteTemplate(w, "choose_action.json", "")
	}
}
