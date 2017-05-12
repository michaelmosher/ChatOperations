package main

import (
	"fmt"
	"io/ioutil"
	"chatOperations/operations"
	"chatOperations/remote"
)

func main() {
	buffer, _ := ioutil.ReadFile("./id_rsa")
	shell, err := remote.NewShell("devuser", buffer)

	var action = operations.Action{
		Id: 1,
		Title: "Deploy",
		Command: "ls -la",
	}

	var server = operations.Server{
		Id: 1,
		Title: "MoB Dev",
		Address: "dev-wp-wise-fs-1.spindance.net",
		Environment: "wp_dev",
	}

	var req = operations.Request{
		Id: 1,
		Requester: "michael.mosher",
		Server: server,
		Action: action,
		Responder: "michael.mosher",
		Approved: true,
	}

	stdout, err := shell.Execute(req)

	fmt.Println(stdout)
	fmt.Println(err)
}
