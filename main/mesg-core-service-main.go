package main

import (
	"os"

	"github.com/nawazish-github/smpl-mesg-service/technology-server"
)

func main() {
	inputs := []string{os.Args[1], os.Args[2], os.Args[3]}

	if inputs[0] != "mesg-core" ||
		inputs[1] != "service" ||
		inputs[2] != "test" {
		return
	}

	server.LaunchServerOn("8080")
}
