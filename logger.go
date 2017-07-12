package main

import (
	"io"
	"time"
)

type Logger struct {
	prefix string // The name of the "thread"
	writer io.Writer
}

// Write is implemented as part of the io.Writer interface
func (l Logger) Write(p []byte) (n int, err error) {
	now := time.Now().UTC()

	// TODO: ljust
	prefix := "[" + l.prefix + "] " +
		now.Format("[15:04:05 Mon 02 Jan UTC] ")

	//TODO Do we need the equivalent of STDOUT.flush in ruby here?
	return l.writer.Write([]byte(prefix + string(p) + "\n"))
}

// This should be used to write strings instead for byte arrays (which is what
// Write methods expects)
func (l Logger) Log(message string) {
	l.Write([]byte(message))
}
