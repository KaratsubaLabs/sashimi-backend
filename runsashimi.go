package main

import (
	"fmt"

	server "github.com/karatsubalabs/sashimi-backend/api"
	"github.com/karatsubalabs/sashimi-backend/sashimi"
)

func main() {

	fmt.Println("Autostarting Sashimi")
	go sashimi.Start()

	server.Start()

}
