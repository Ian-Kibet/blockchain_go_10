package main

import (
	"fmt"
)

func (cli *CLI) startHttpServer(port string) {
	fmt.Printf("Starting http server on port: %s\n", port)

	RunHttpServer(port)
	// StartServer(nodeID, minerAddress)
}
