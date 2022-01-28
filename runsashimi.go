package main

import (
	server "github.com/karatsubalabs/sashimi-backend/api"
	"github.com/karatsubalabs/sashimi-backend/sashimi"
)

func main() {
	go sashimi.Start()
	server.Start()
}
