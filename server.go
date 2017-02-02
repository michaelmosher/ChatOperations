package main

import (
	"fmt"
	"log"
	"net/http"
	"slackApi/api"
	"os"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome!\n")
}

func main() {
	port := os.Getenv("PORT")
	mux := http.NewServeMux()

	mux.HandleFunc("/", Index)
	mux.HandleFunc("/api/hello", slackApi.Hello)
	mux.HandleFunc("/api/goodbye", slackApi.Goodbye)
	mux.HandleFunc("/operations", slackApi.Operations)

	log.Fatal(http.ListenAndServe(":" + port, mux))
}
