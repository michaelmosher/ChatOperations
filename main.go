package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"chatoperations/api"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome!\n")
}

func main() {
	port := os.Getenv("PORT")
	mux := http.NewServeMux()

	cfg := slackApi.ApiConfig{
		VerificationToken: os.Getenv("VerificationToken"),
		WebhookUrl:        os.Getenv("WebhookUrl"),
		DatabaseUrl:       os.Getenv("DATABASE_URL"),
		NetClientTimeout:  time.Second * 10,
	}

	api := slackApi.New(cfg)

	mux.HandleFunc("/", Index)
	mux.HandleFunc("/operations", api.Operations)

	log.Fatal(http.ListenAndServe(":"+port, mux))
}
