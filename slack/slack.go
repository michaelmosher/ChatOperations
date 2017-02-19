package slack

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"time"

	"chatoperations/operations"
)

type Config struct {
	TemplatesGlob string
	WebhookUrl    string
}

type Slack struct {
	WebhookUrl    string
	NetClient     *http.Client
	JsonTemplates *template.Template
}

func New(cfg Config) *Slack {
	netClient := &http.Client{
		Timeout: time.Second * 10,
	}
	// load templates
	templates := template.Must(template.ParseGlob(cfg.TemplatesGlob))

	return &Slack{
		WebhookUrl:    cfg.WebhookUrl,
		NetClient:     netClient,
		JsonTemplates: templates,
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

func (slack *Slack) NotifyRequestSubmitted(o operations.Request) error {
	return slack.httpPost(slack.WebhookUrl, "ops_request_submitted.json", o)
}

func (slack *Slack) NotifyRequestApproved(o operations.Request) error {
	return slack.httpPost(o.Response_url, "request_approved.json", o)
}

func (slack *Slack) NotifyRequestRejected(o operations.Request) error {
	return slack.httpPost(o.Response_url, "request_rejected.json", o)
}

func (slack *Slack) NotifyError(url string, originalErr error) {
	newError := slack.httpPost(url, "something_went_wrong.json", originalErr.Error)

	if newError != nil {
		fmt.Printf("Failed to send error notification: %s\nOriginal error: %s",
			newError, originalErr)
	}
}
