package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"chatoperations/application"
	"chatoperations/database"
	"chatoperations/slack"
	"chatoperations/web"
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome!\n")
}

func main() {
	port := os.Getenv("PORT")
	mux := http.NewServeMux()

	state, err := database.New(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	notifier := slack.New(slack.Config{
		TemplatesGlob: "slack/templates/*.json",
		WebhookUrl:    os.Getenv("WebhookUrl"),
	})

	app := application.OperationsInteractor{
		ActionStore:  state.NewActionRepo(),
		ServerStore:  state.NewServerRepo(),
		RequestStore: state.NewRequestRepo(),
		Notifier:     notifier,
	}

	api := web.Server{
		VerificationToken: os.Getenv("VerificationToken"),
		OpsInteractor:     app,
	}

	mux.HandleFunc("/", index)
	mux.HandleFunc("/operations", api.Operations)

	log.Fatal(http.ListenAndServe(":"+port, mux))
}
