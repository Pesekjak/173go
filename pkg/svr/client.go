package svr

import (
	"fmt"

	"github.com/Pesekjak/173go/pkg/base"
	"github.com/Pesekjak/173go/pkg/log"
	"github.com/Pesekjak/173go/pkg/net"
	"github.com/Pesekjak/173go/pkg/prot"
	"github.com/Pesekjak/173go/pkg/world"
	"github.com/Pesekjak/173go/pkg/world/material"
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

	err := c.connection.WritePacket(&prot.PacketOutLogin{
		EntityId:   base.NextEntityId(),
		ServerName: "", // empty on Notchian
		MapSeed:    0,  // unused by client
		Dimension:  byte(defaultWorld.Dimension),
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
		Time: defaultWorld.Time,
	}, false)
	if err != nil {
		return err
	}

	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			chunk := world.NewChunk(world.NewChunkPos(int32(i), int32(j)), 127, func(blocks []world.Block) error {
				return nil
			})
			for x := 0; x < 16; x++ {
				for z := 0; z < 16; z++ {
					block, _ := chunk.GetBlock(int32(x), 0, int32(z))
					block.Set(material.Stone, 0)
				}
			}
			chunk.Load(c.connection)
		}
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
