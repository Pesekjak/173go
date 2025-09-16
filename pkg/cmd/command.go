package cmd

import (
	"strings"
	"unicode/utf8"

	"github.com/Pesekjak/173go/pkg/log"
)

// CommandManager manager of server commands.
type CommandManager struct {
	commands map[string]Command
	logger   *log.Logger
}

// Command executable command by any CommandSender
type Command struct {
	// Label used to run the command
	Label string
	// Usage string displayed if the command execution fails
	Usage string
	// Permission required to execute the command
	Permission string
	// Handler for the command logic. Source of the command and arguments split at space and trimmed are provided
	Handler func(sender CommandSender, args []string) bool
}

// NewCommandManager provides new empty CommandManager
func NewCommandManager(logger *log.Logger) *CommandManager {
	return &CommandManager{commands: make(map[string]Command), logger: logger}
}

// RegisterCommand registers new command, fails and returns false if there is already a command with the same
// label registered, else returns true
func (cm *CommandManager) RegisterCommand(cmd Command) bool {
	if _, ok := cm.commands[cmd.Label]; ok {
		return false
	}
	cm.commands[cmd.Label] = cmd
	return true
}

// ExecuteCommand executes the command with given source and buffer.
// Buffer is a whole command string (without the starting '/').
// Returns true if the command execution was successful, else false.
// Sends extra messages to the source in case the execution fails (about missing permissions, incorrect usage...)
func (cm *CommandManager) ExecuteCommand(sender CommandSender, buffer string) bool {
	buffer = strings.TrimSpace(buffer)
	if utf8.RuneCountInString(buffer) == 0 {
		return false // empty command
	}

	args := strings.Split(buffer, " ")
	for i, arg := range args {
		args[i] = strings.TrimSpace(arg)
	}

	label := args[0] // always present
	var parameters []string
	if len(args) == 1 {
		parameters = []string{}
	} else {
		parameters = args[1:]
	}

	cmd, ok := cm.commands[label]
	if !ok {
		sender.SendMessage("unknown command: ", label)
		return false
	}

	if !sender.HasPermission(cmd.Permission) {
		sender.SendMessage("you do not have permissions to execute this command")
	}

	cm.logger.Info(sender, " executed command: /", buffer)
	result := cmd.Handler(sender, parameters)

	if !result {
		sender.SendMessage("incorrect usage: " + cmd.Usage)
		return false
	}

	return true
}
