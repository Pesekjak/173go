package svr

import "github.com/Pesekjak/173go/pkg/cmd"

func registerCommands(server *Server) {
	server.CommandManager.RegisterCommand(cmd.Command{
		Label:      "stop",
		Usage:      "/stop",
		Permission: "server.stop",
		Handler: func(sender cmd.CommandSender, args []string) bool {
			if len(args) != 0 {
				return false
			}
			server.Stop()
			return true
		},
	})
}
