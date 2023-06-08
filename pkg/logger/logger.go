package logger

import (
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

type Level int8

const (
	LevelInfo Level = iota
	LevelDebug
	LevelError
	LevelFatal
	LevelOff
)

func (l Level) String() string {
	switch l {
	case LevelInfo:
		return "INFO"
	case LevelDebug:
		return "DEBUG"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return ""
	}
}

type Logger struct {
	out      io.Writer
	minLevel Level
	mu       sync.Mutex
}

func NewLogger(out io.Writer, minLevel Level) *Logger {
	return &Logger{
		out:      out,
		minLevel: minLevel,
	}
}

func (l *Logger) PrintInfo(message string) {
	l.print(LevelInfo, message)
}

func (l *Logger) PrintDebug(message string) {
	l.print(LevelDebug, message)
}

func (l *Logger) PrintError(err error) {
	l.print(LevelError, err.Error())
}

func (l *Logger) PrintFatal(err error) {
	l.print(LevelFatal, err.Error())
	os.Exit(1)
}

func (l *Logger) print(level Level, message string) (int, error) {
	if level < l.minLevel {
		return 0, nil
	}

	aux := struct {
		Level      string            `json:"level"`
		Time       string            `json:"time"`
		Message    string            `json:"message"`
		Properties map[string]string `json:"properties,omitempty"`
		Trace      string            `json:"trace,omitempty"`
	}{
		Level: level.String(),

		Message: message,
	}

	time1 := time.Now().UTC().Format(time.RFC3339)

	if level >= LevelError {
		aux.Trace = string(debug.Stack())
	}

	outMessage := fmt.Sprintf("%s: [%v]: %v", level.String(), time1, message)

	line := []byte(outMessage)

	return l.out.Write(append(line, '\n'))
}
