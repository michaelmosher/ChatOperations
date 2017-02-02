package slackApi

import (
	"fmt"
	"net/http"
)

func MyLibHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "This is where my API will go!")
}

func Hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "hello!")
}

func Goodbye(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "good bye!")
}

func Operations(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "user: " + r.PostFormValue("user_name") + "<br>" +
		   "text: " + r.PostFormValue("text"))
}
