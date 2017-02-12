package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"slackApi/api"
	"time"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome!\n")
}

func main() {
	port := os.Getenv("PORT")
	mux := http.NewServeMux()

	db, _ := slackApi.NewDB(os.Getenv("DATABASE_URL"))
	netClient := &http.Client{
		Timeout: time.Second * 10,
	}

	env := slackApi.Env{
		VerificationToken: os.Getenv("VerificationToken"),
		WebhookUrl:        os.Getenv("WebhookUrl"),
		Db:                db,
		NetClient:         netClient,
	}

	mux.HandleFunc("/", Index)
	mux.HandleFunc("/operations", env.Operations)

	log.Fatal(http.ListenAndServe(":"+port, mux))
}
