package main

import (
	"log"
	"net/http"

	agent "github.com/peterdeme/noobernetes/agent/routes"
)

func main() {
	http.HandleFunc("containers/start", agent.StartContainers)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
