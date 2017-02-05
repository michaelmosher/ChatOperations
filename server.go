package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"slackApi/api"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome!\n")
}

func main() {
	port := os.Getenv("PORT")
	mux := http.NewServeMux()

	mux.HandleFunc("/", Index)
	mux.HandleFunc("/operations", slackApi.Operations)

	// /oauth to start handshake
	// /operations to do stuff

	log.Fatal(http.ListenAndServe(":"+port, mux))
}
