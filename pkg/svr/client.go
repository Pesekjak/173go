package svr

import (
	"fmt"

	"github.com/Pesekjak/173go/pkg/base"
	"github.com/Pesekjak/173go/pkg/log"
	"github.com/Pesekjak/173go/pkg/net"
	"github.com/Pesekjak/173go/pkg/prot"
	"github.com/Pesekjak/173go/pkg/world"
	"github.com/Pesekjak/173go/pkg/world/entity_data"
)

type Client struct {
	server     *Server
	connection *net.Connection

	logger *log.Logger

	id       int32
	username string
	location world.Location
	world    *world.World
}

func NewClient(server *Server, connection *net.Connection) *Client {
	return &Client{
		server:     server,
		connection: connection,

		logger: server.Console.ChildLogger("client"),
	}
}

func (c *Client) OnLogin(packet *prot.PacketInLogin) error {
	if packet.Username != c.username {
		return fmt.Errorf("client %v tried to login with username '%v'", c, packet.Username)
	}
	if packet.Protocol != 14 {
		return fmt.Errorf("unsupported protocol version: %v", packet.Protocol)
	}

	defaultWorld := c.server.defaultWorld

	c.id = base.NextEntityId()
	c.location = world.NewLocation(0, 0, 0, 0, 0)
	c.world = defaultWorld

	err := c.connection.WritePacket(&prot.PacketOutLogin{
		EntityId:   c.id,
		ServerName: "", // empty on Notchian
		MapSeed:    0,  // unused by client
		Dimension:  byte(defaultWorld.Dimension()),
	}, false)
	if err != nil {
		return err
	}

	spawnPoint := defaultWorld.SpawnPoint

	err = c.connection.WritePacket(&prot.PacketOutSpawnPosition{
		X: spawnPoint.X,
		Y: spawnPoint.Y,
		Z: spawnPoint.Z,
	}, false)
	if err != nil {
		return err
	}

	err = c.connection.WritePacket(&prot.PacketOutTimeUpdate{
		Time: defaultWorld.Time(),
	}, false)
	if err != nil {
		return err
	}

	if err = defaultWorld.SpawnPlayer(c); err != nil {
		return err
	}

	return c.connection.WritePacket(&prot.PacketOutPlayerPositionAndLook{
		X:        float64(spawnPoint.X),
		Stance:   67.240000009536743,
		Y:        float64(spawnPoint.Y),
		Z:        float64(spawnPoint.Z),
		Yaw:      0,
		Pitch:    0,
		OnGround: false,
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

func (c *Client) OnPlayerGround(*prot.PacketInPlayerGround) error {
	return nil
}

func (c *Client) OnPlayerPosition(*prot.PacketInPlayerPosition) error {
	return nil
}

func (c *Client) OnPlayerLook(*prot.PacketInPlayerLook) error {
	return nil
}

func (c *Client) OnPlayerPositionAndLook(*prot.PacketInPlayerPositionAndLook) error {
	return nil
}

func (c *Client) Id() int32 {
	return c.id
}

func (c *Client) Location() world.Location {
	return c.location
}

func (c *Client) World() *world.World {
	return c.world
}

func (c *Client) EntityType() world.EntityType {
	return world.Player
}

func (c *Client) IsAlive() bool {
	return true
}

func (c *Client) Health() uint32 {
	return 20
}

func (c *Client) Metadata() entity_data.Metadata {
	return struct{}{}
}

func (c *Client) Username() string {
	return c.username
}

func (c *Client) IsOnline() bool {
	return c.connection.IsActive()
}

func (c *Client) Connection() *net.Connection {
	return c.connection
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
