package main

import (
	"fmt"
	"log"
	"net/http"
	"slackApi/api"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome!\n")
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", Index)
	mux.HandleFunc("/api/hello", slackApi.Hello)
	mux.HandleFunc("/api/goodbye", slackApi.Goodbye)

	log.Fatal(http.ListenAndServe(":8080", mux))
}
