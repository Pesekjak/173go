package system

type Command int

type Message struct {
	Command
	Reason interface{}
}

const (
	Stop Command = iota
	Fail
)

func Make(command Command, reason interface{}) Message {
	return Message{
		Command: command,
		Reason:  reason,
	}
}
