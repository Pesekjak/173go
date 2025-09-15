package net

import (
	"fmt"
	"net"
	"strconv"

	"github.com/Pesekjak/173go/pkg/log"
	"github.com/Pesekjak/173go/pkg/prot"
	"github.com/Pesekjak/173go/pkg/system"
)

type Network struct {
	host string
	port int

	logger *log.Logger

	report chan system.Message
}

func NewNetwork(host string, port int, logger *log.Logger, report chan system.Message) *Network {
	return &Network{
		host:   host,
		port:   port,
		logger: logger,
		report: report,
	}
}

func (n *Network) Start(handlerSupply func(conn *Connection) (prot.PacketHandler, error)) error {
	address := n.host + ":" + strconv.Itoa(n.port)
	ser, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return fmt.Errorf("address resolution failed: %v", err)
	}

	tcp, err := net.ListenTCP("tcp", ser)
	if err != nil {
		return fmt.Errorf("failed to bind: %v", err)
	}

	go func() {
		n.logger.Info("accepting connections on ", address)
		for {
			clientTcp, err := tcp.AcceptTCP()

			if err != nil {
				break
			}

			_ = clientTcp.SetNoDelay(true)
			_ = clientTcp.SetKeepAlive(true)

			conn := NewConnection(clientTcp, n.logger, n.report)
			handler, err := handlerSupply(conn)
			if err != nil {
				n.logger.Severe("failed to accept connection: ", tcp)
				return
			}
			go func() {
				// more accurate information is logged within the connection itself
				_ = conn.startListening(handler)
			}()
		}
	}()

	return nil
}
