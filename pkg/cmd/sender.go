package cmd

// CommandSender can execute server commands.
type CommandSender interface {
	// Name of the sender
	Name() string
	// SendMessage sends a message to the sender.
	SendMessage(message ...interface{})
	// HasPermission checks if the sender has a given permission.
	HasPermission(perm string) bool
}
