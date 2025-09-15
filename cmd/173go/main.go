package main

import (
	"github.com/Pesekjak/173go/pkg/svr"
)

func main() {
	server := svr.NewServer()
	server.Start()
}
