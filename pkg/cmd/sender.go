package cmd

type Sender interface {
	Name() string
	SendMessage(message ...interface{})
	HasPermission(perm string) bool
}
