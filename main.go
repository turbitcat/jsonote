package main

import (
	"log"

	"github.com/turbitcat/jsonote/v2/api"
)

func main() {
	server := api.New()
	log.Fatalln(server.Run())
}
