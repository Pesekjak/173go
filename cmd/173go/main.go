package main

import (
	"fmt"

	"github.com/Pesekjak/173go/pkg/svr"
)

func main() {
	server, err := svr.NewServer()
	if err != nil {
		fmt.Printf("Failed to create the server: %v\n", err)
		return
	}
	server.Start()
}
