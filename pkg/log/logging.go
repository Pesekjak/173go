package log

import (
	"fmt"
	"io"
	"time"

	"github.com/Pesekjak/173go/pkg/base"
	"github.com/fatih/color"
)

type Level int

const (
	Info Level = iota
	Warn
	Severe
	Debug
)

var BasicLevels = []Level{Info, Warn, Severe}
var AllLevels = []Level{Info, Warn, Severe, Debug}

var levelMetadata = map[Level]struct {
	Str   string
	Color *color.Color
}{
	Info:   {"INFO", color.New(color.FgCyan)},
	Warn:   {"WARN", color.New(color.FgYellow)},
	Severe: {"ERROR", color.New(color.FgRed)},
	Debug:  {"DEBUG", color.New(color.FgMagenta)},
}

type Logger struct {
	name          string
	writer        io.Writer
	enabledLevels map[Level]bool
}

func NewLogger(name string, writer io.Writer, levels ...Level) *Logger {
	enabled := make(map[Level]bool, len(levels))
	for _, lvl := range levels {
		enabled[lvl] = true
	}
	return &Logger{
		name:          name,
		writer:        writer,
		enabledLevels: enabled,
	}
}

func (l *Logger) Name() string {
	return l.name
}

func (l *Logger) Levels() []Level {
	levels := make([]Level, 0, len(l.enabledLevels))
	for lvl := range l.enabledLevels {
		levels = append(levels, lvl)
	}
	return levels
}

func (l *Logger) logf(level Level, format string, a ...interface{}) {
	if !l.enabledLevels[level] {
		return
	}

	meta, ok := levelMetadata[level]
	if !ok {
		meta = struct {
			Str   string
			Color *color.Color
		}{"UNKWN", color.New(color.FgWhite)}
	}

	_, _ = fmt.Fprintf(l.writer, "[%s] [%s] [%s] %s\n",
		color.HiGreenString(currentTimeAsText()),
		meta.Color.Sprint(meta.Str),
		color.WhiteString(l.name),
		fmt.Sprintf(format, a...),
	)
}

func (l *Logger) Info(message ...interface{}) {
	l.logf(Info, "%s", base.ConvertToString(message...))
}

func (l *Logger) Warn(message ...interface{}) {
	l.logf(Warn, "%s", base.ConvertToString(message...))
}

func (l *Logger) Severe(message ...interface{}) {
	l.logf(Severe, "%s", base.ConvertToString(message...))
}

func (l *Logger) Debug(message ...interface{}) {
	l.logf(Debug, "%s", base.ConvertToString(message...))
}

func (l *Logger) InfoF(format string, a ...interface{}) {
	l.logf(Info, format, a...)
}

func (l *Logger) WarnF(format string, a ...interface{}) {
	l.logf(Warn, format, a...)
}

func (l *Logger) SevereF(format string, a ...interface{}) {
	l.logf(Severe, format, a...)
}

func (l *Logger) DebugF(format string, a ...interface{}) {
	l.logf(Debug, format, a...)
}

func currentTimeAsText() string {
	return time.Now().Format("15:04:05")
}
