package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"chatoperations/database"
	"chatoperations/server"
	"chatoperations/slack"
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

	api := server.New(server.Config{
		TemplatesGlob:     "server/templates/*.json",
		VerificationToken: os.Getenv("VerificationToken"),
		State:             state,
		Notifier:          notifier,
	})

	mux.HandleFunc("/", index)
	mux.HandleFunc("/operations", api.Operations)

	log.Fatal(http.ListenAndServe(":"+port, mux))
}
