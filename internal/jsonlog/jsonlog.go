package jsonlog

import (
	"encoding/json"
	"fmt"
	"io"
	"runtime/debug"
	"sync"
	"time"
)

type Level int8

const (
	LevelInfo Level = iota
	LevelError
	LevelFatal
	LevelOff
)

func (l Level) String() string {
	switch l {
	case LevelInfo:
		return "INFO"
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

// New creates a new logger
//
// Parameters:
//   - out: io.Writer
//   - minLevel: minimum Level of the information
//
// Returns
//   - Pointer to the newly created Logger object
//
// This method is used for creating a new Logger object
func New(out io.Writer, minLevel Level) *Logger {
	return &Logger{
		out:      out,
		minLevel: minLevel,
	}
}

func (l *Logger) print(level Level, message string, properties map[string]string) (int, error) {
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
		Level:      level.String(),
		Time:       time.Now().UTC().Format(time.RFC3339),
		Message:    message,
		Properties: properties,
	}
	if level >= LevelError {
		aux.Trace = string(debug.Stack())
	}
	var line []byte
	line, err := json.Marshal(aux)
	if err != nil {
		line = []byte(LevelError.String() + ": unable to marshal log message" + err.Error())
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.out.Write(append(line, '\n'))
}

// Write writes a log message at the error level.
//
// Parameters:
// - message: A byte slice containing the log message.
//
// Returns:
//   - n: The number of bytes written, which is typically the length of the message.
//   - err: An error if the write operation fails, or nil if successful.
//
// This method is typically used to satisfy the `io.Writer` interface, allowing the logger
// to be used in contexts where an `io.Writer` is required.
func (l *Logger) Write(message []byte) (n int, err error) {
	return l.print(LevelError, string(message), nil)
}

// PrintInfo logs an informational message with additional properties.
//
// Parameters:
//   - message: A string containing the informational message to be logged.
//   - properties: A map of key-value pairs containing additional properties to be included with the log message.
//
// This method logs the message at the info level and can include extra context through the properties map.
// The properties map allows for attaching relevant data to the log entry, such as user IDs or request details.
func (l *Logger) PrintInfo(message string, properties map[string]string) {
	_, _ = l.print(LevelInfo, fmt.Sprint(message), properties)
}

// PrintError logs an error message with additional properties.
//
// Parameters:
//   - err: The error to be logged. The error message will be extracted and logged.
//   - properties: A map of key-value pairs containing additional properties to be included with the log message.
//
// This method logs the error message at the error level and can include extra context through the properties map.
// The properties map allows for attaching relevant data to the log entry, such as error codes or request details.
func (l *Logger) PrintError(err error, properties map[string]string) {
	_, _ = l.print(LevelError, err.Error(), properties)
}

// PrintFatal logs a fatal error message with additional properties.
//
// Parameters:
//   - err: The error to be logged. The error message will be extracted and logged as a fatal error.
//   - properties: A map of key-value pairs containing additional properties to be included with the log message.
//
// This method logs the error message at the fatal level, indicating a critical issue that typically requires immediate attention.
// The properties map allows for attaching relevant data to the log entry, such as error codes or request details.
func (l *Logger) PrintFatal(err error, properties map[string]string) {
	_, _ = l.print(LevelFatal, err.Error(), properties)
}
