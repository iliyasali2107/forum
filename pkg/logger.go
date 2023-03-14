package pkg

import (
	"fmt"
	"log"
	"time"
)

type Logger struct {
	prefix string
}

// NewLogger creates a new Logger with a prefix
func NewLogger(prefix string) *Logger {
	return &Logger{prefix: prefix}
}

// Info logs an info message
func (l *Logger) Info(format string, v ...interface{}) {
	l.log("[INFO]", format, v...)
}

// Error logs an error message
func (l *Logger) Error(format string, v ...interface{}) {
	l.log("[ERROR]", format, v...)
}

// log logs a message with a specific prefix
func (l *Logger) log(prefix, format string, v ...interface{}) {
	log.Printf("%s %s %s: %s\n", l.prefix, time.Now().Format(time.RFC3339), prefix, fmt.Sprintf(format, v...))
}
