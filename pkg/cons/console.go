package cons

import (
	"bufio"
	"io"
	"strings"
	"sync"

	"github.com/Pesekjak/173go/pkg/base"
	"github.com/Pesekjak/173go/pkg/chat"
	"github.com/Pesekjak/173go/pkg/log"
	"github.com/fatih/color"
)

// Console represents a console command sender
type Console struct {
	reader io.Reader
	writer io.Writer

	// IChannel for the console input.
	// Can be used to execute console commands from outside
	IChannel chan string

	closer sync.Once

	log.Logger
}

// NewConsole creates new console wrapped around given reader and writer, logging messages at provided levels.
func NewConsole(reader io.Reader, writer io.Writer, levels ...log.Level) *Console {
	console := &Console{
		reader:   reader,
		writer:   writer,
		IChannel: make(chan string),
		Logger:   *log.NewLogger("server", writer, levels...),
	}
	return console
}

// Start starts the console using the given command handler.
// If any data is sent to the IChannel, it is piped to the given handler.
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

// Stop closes the console IChannel.
func (c *Console) Stop() {
	c.closer.Do(func() {
		close(c.IChannel)
	})
}

func (c *Console) Name() string {
	return "server"
}

func (c *Console) SendMessage(message ...interface{}) {
	c.Info(toTerminal(base.ConvertToString(message...))) // for send message we translate colored messages
}

func (c *Console) HasPermission(string) bool {
	return true // console has all permissions
}

func (c *Console) String() string {
	return c.Name()
}

var (
	colorMap = map[rune]color.Attribute{
		'0': color.FgBlack,
		'1': color.FgBlue,
		'2': color.FgGreen,
		'3': color.FgCyan,
		'4': color.FgRed,
		'5': color.FgMagenta,
		'6': color.FgYellow,
		'7': color.FgWhite,
		'8': color.FgHiBlack,
		'9': color.FgHiBlue,
		'a': color.FgHiGreen,
		'b': color.FgHiCyan,
		'c': color.FgHiRed,
		'd': color.FgHiMagenta,
		'e': color.FgHiYellow,
		'f': color.FgHiWhite,
	}

	styleMap = map[rune]color.Attribute{
		'l': color.Bold,
		'n': color.Underline,
		'o': color.Italic,
		'm': color.CrossedOut,
	}
)

func toTerminal(m string) string {
	var resultBuilder strings.Builder
	var textSegmentBuilder strings.Builder
	var currentAttributes []color.Attribute

	flush := func() {
		if textSegmentBuilder.Len() > 0 {
			c := color.New(currentAttributes...)
			resultBuilder.WriteString(c.Sprint(textSegmentBuilder.String()))
			textSegmentBuilder.Reset()
		}
	}

	runes := []rune(m)
	for i := 0; i < len(runes); i++ {
		if runes[i] == chat.ColorSymbol && i+1 < len(runes) {
			flush()
			code := runes[i+1]
			i++

			if newColor, isColor := colorMap[code]; isColor {
				currentAttributes = []color.Attribute{newColor}
			} else if newStyle, isStyle := styleMap[code]; isStyle {
				currentAttributes = append(currentAttributes, newStyle)
			} else if code == 'r' {
				currentAttributes = []color.Attribute{}
			}
			// §k (magic) is intentionally ignored
			continue
		}

		textSegmentBuilder.WriteRune(runes[i])
	}

	flush()
	return resultBuilder.String()
}
