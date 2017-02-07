package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"database/sql"
	"slackApi/api"
)

type Env struct {
	db *sql.DB
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome!\n")
}

func main() {
	// db, err := slackApi.NewDB(os.Getenv("DATABASE_URL"))
	// if err != nil {
	// 	log.Panic(err)
	// }
	// env := &Env{db: db}

	port := os.Getenv("PORT")
	mux := http.NewServeMux()

	mux.HandleFunc("/", Index)
	mux.HandleFunc("/operations", slackApi.Operations)

	log.Fatal(http.ListenAndServe(":" + port, mux))
}
