package slack

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"time"
)

type Config struct {
	WebhookUrl	string
	NetClientTimeout  time.Duration
}

type Slack struct {
	WebhookUrl        string
	NetClient         *http.Client
	JsonTemplates     *template.Template
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

func New(cfg Config) Slack {
	netClient := &http.Client{
		Timeout: cfg.NetClientTimeout,
	}
	// load templates
	templates := template.Must(template.ParseGlob("templates/*.json"))

	return Slack{
		WebhookUrl:        cfg.WebhookUrl,
		NetClient:         netClient,
		JsonTemplates:     templates,
	}
}

func (slack *Slack) httpPost(url string, template_name string, data interface{}) error {
	reader, writer := io.Pipe()

	// writing without a reader will deadlock so write in a goroutine
	go func() {
		defer writer.Close()
		slack.JsonTemplates.ExecuteTemplate(writer, template_name, data)
	}()

	resp, err := slack.NetClient.Post(url, "application/json", reader)

	if err != nil {
		return err
	}

	defer resp.Body.Close()
	return nil
}

func (slack *Slack) NotifyRequestSubmitted(opsRequest OperationsRequest) error {
	return slack.httpPost(slack.WebhookUrl, "ops_request_submitted.json", opsRequest)
}

func (slack *Slack) NotifyRequestApproved(opsRequest OperationsRequest) error {
	return slack.httpPost(opsRequest.Response_url, "request_approved.json", opsRequest)
}

func (slack *Slack) NotifyRequestRejected(opsRequest OperationsRequest) error {
	return slack.httpPost(opsRequest.Response_url, "request_rejected.json", opsRequest)
}

func (slack *Slack) NotifyError(url string, originalErr error) {
	newError := slack.httpPost(url, "something_went_wrong.json", originalErr)

	if newError != nil {
		fmt.Printf("Failed to send error notification: %s\nOriginal error: %s",
					newError, originalErr)
	}
}
