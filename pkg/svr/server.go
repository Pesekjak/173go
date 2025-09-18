package svr

import (
	"os"

	"github.com/Pesekjak/173go/pkg/cmd"
	"github.com/Pesekjak/173go/pkg/cons"
	"github.com/Pesekjak/173go/pkg/log"
	"github.com/Pesekjak/173go/pkg/net"
	"github.com/Pesekjak/173go/pkg/prot"
	"github.com/Pesekjak/173go/pkg/system"
	"github.com/Pesekjak/173go/pkg/world"
)

type Server struct {
	message chan system.Message

	Config

	*cons.Console
	*cmd.CommandManager

	defaultWorld *world.World

	clients []Client
}

func NewServer() (*Server, error) {
	message := make(chan system.Message)

	config := NewDefaultConfig()
	console := cons.NewConsole(os.Stdin, os.Stdout, log.BasicLevels...)

	defaultWorld, err := world.NewWorld()
	if err != nil {
		return nil, err
	}

	commandManager := cmd.NewCommandManager(console.ChildLogger("cmd"))

	server := &Server{
		message: message,

		Config: config,

		Console:        console,
		CommandManager: commandManager,

		defaultWorld: defaultWorld,

		clients: make([]Client, 0, 8),
	}

	return server, nil
}

func (s *Server) Start() {
	s.Console.Start(func(cmd string) {
		s.CommandManager.ExecuteCommand(s.Console, cmd)
	})
	s.Console.Info("starting 173go server...")

	registerCommands(s)

	network := net.NewNetwork(s.Config.Address, s.Config.Port, s.Console.ChildLogger("network"), s.message)
	if err := network.Start(func(conn *net.Connection) (prot.PacketHandler, error) {
		return NewClient(s, conn), nil
	}); err != nil {
		s.Console.Severe("failed to start the network server: ", err)
		return
	}

	for x := int32(-4); x <= 4; x++ {
		for z := int32(-4); z <= 4; z++ {
			if _, err := s.defaultWorld.LoadChunk(world.NewChunkPos(x, z)); err != nil {
				s.Console.Severe("failed to load the chunk at ", x, ";", z, ": ", err)
			}
		}
	}

	s.Console.Info("server is running")
	s.wait()
}

func (s *Server) Stop() {
	s.message <- system.Make(system.Stop, nil)
}

func (s *Server) wait() {
	for {
		select {
		case command := <-s.message:
			switch command.Command {
			case system.Stop:
				s.Console.Info("stopping server")
				s.terminate()
				return
			case system.Fail:
				s.Console.Severe("internal server error: ", command.Reason)
				s.Console.Info("stopping server")
				s.terminate()
				return
			}
		}
	}
}

func (s *Server) terminate() {
	s.Console.Stop()
}
