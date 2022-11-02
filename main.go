package main

import "github.com/turbitcat/jsonote/v2/api"

func main() {
	server := api.New()
	server.Run()
}
