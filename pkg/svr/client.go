package svr

import (
	"fmt"

	"github.com/Pesekjak/173go/pkg/base"
	"github.com/Pesekjak/173go/pkg/log"
	"github.com/Pesekjak/173go/pkg/net"
	"github.com/Pesekjak/173go/pkg/prot"
	"github.com/Pesekjak/173go/pkg/world"
)

type Client struct {
	server     *Server
	connection *net.Connection

	logger *log.Logger

	username string
}

func NewClient(server *Server, connection *net.Connection) *Client {
	return &Client{
		server:     server,
		connection: connection,

		logger: server.Console.ChildLogger("client_handler"),
	}
}

func (c *Client) OnLogin(packet *prot.PacketInLogin) error {
	if packet.Username != c.username {
		return fmt.Errorf("client %v tried to login with username '%v'", c, packet.Username)
	}
	if packet.Protocol != 14 {
		return fmt.Errorf("unsupported protocol version: %v", packet.Protocol)
	}
	return c.connection.WritePacket(&prot.PacketOutLogin{
		EntityId:   base.NextEntityId(),
		ServerName: "", // empty on Notchian
		MapSeed:    0,  // unused by client
		Dimension:  world.Overworld,
	}, true)
}

func (c *Client) OnKeepAlive(*prot.PacketInKeepAlive) error {
	return c.connection.WritePacket(&prot.PacketOutKeepAlive{}, true)
}

func (c *Client) OnHandShake(packet *prot.PacketInHandShake) error {
	c.logger.Info("new handshake with: '", packet.Username, "'")
	c.username = packet.Username
	return c.connection.WritePacket(&prot.PacketOutHandShake{Hash: "-"}, true)
}

func (c *Client) OnPlayerPosition(*prot.PacketInPlayerPosition) error {
	return nil
}

func (c *Client) OnPlayerPositionAndLook(*prot.PacketInPlayerPositionAndLook) error {
	return nil
}

func (c *Client) Disconnect(err error) {
	c.connection.Close(err)
}

func (c *Client) Kick(reason string) {
	c.connection.CloseWith(nil, reason)
}

func (c *Client) String() string {
	if c.username == "" {
		return "unknown client?"
	}
	return c.username
}
