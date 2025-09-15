package cons

import (
	"bufio"
	"io"
	"sync"

	"github.com/Pesekjak/173go/pkg/base"
	"github.com/Pesekjak/173go/pkg/log"
)

type Console struct {
	reader io.Reader
	writer io.Writer

	IChannel chan string

	closer sync.Once

	log.Logger
}

func NewConsole(reader io.Reader, writer io.Writer, levels ...log.Level) *Console {
	console := &Console{
		reader:   reader,
		writer:   writer,
		IChannel: make(chan string),
		Logger:   *log.NewLogger("server", writer, levels...),
	}
	return console
}

func (c *Console) Start(commandHandler func(string)) {
	go func() {
		scanner := bufio.NewScanner(c.reader)
		for scanner.Scan() {
			c.IChannel <- scanner.Text()
		}
	}()
	go func() {
		for cmd := range c.IChannel {
			commandHandler(cmd)
		}
	}()
}

func (c *Console) Stop() {
	c.closer.Do(func() {
		close(c.IChannel)
	})
}

func (c *Console) ChildLogger(name string) *log.Logger {
	return log.NewLogger(name, c.writer, c.Logger.Levels()...)
}

func (c *Console) Name() string {
	return "server"
}

func (c *Console) SendMessage(message ...interface{}) {
	c.Info(base.ConvertToString(message...))
}

func (c *Console) HasPermission(string) bool {
	return true
}
