package cmd

import (
	"strings"
	"unicode/utf8"
)

type CommandManager struct {
	commands map[string]Command
}

type Command struct {
	Label      string
	Usage      string
	Permission string
	Handler    func(sender Sender, args []string) bool
}

func NewCommandManager() *CommandManager {
	return &CommandManager{commands: make(map[string]Command)}
}

func (cm *CommandManager) RegisterCommand(cmd Command) bool {
	if _, ok := cm.commands[cmd.Label]; ok {
		return false
	}
	cm.commands[cmd.Label] = cmd
	return true
}

func (cm *CommandManager) ExecuteCommand(sender Sender, buffer string) bool {
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

	result := cmd.Handler(sender, parameters)

	if !result {
		sender.SendMessage("incorrect usage: " + cmd.Usage)
		return false
	}

	return true
}
