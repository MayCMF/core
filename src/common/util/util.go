package util

import (
	"fmt"
	"os"
	"time"
)

var (
	pid = os.Getpid()
)

// NewTraceID - Create tracking ID
func NewTraceID() string {
	return fmt.Sprintf("trace-id-%d-%s",
		pid,
		time.Now().Format("2006.01.02.15.04.05.999999"))
}
