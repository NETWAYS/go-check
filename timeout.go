package check

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var timeoutEnabled bool

// Start the timeout and signal handler in a goroutine
func (c *Config) EnableTimeoutHandler() {
	go HandleTimeout(c.Timeout)
}

// Helper for a goroutine, to wait for signals and timeout, and exit with a proper code
func HandleTimeout(timeout int) {
	if timeoutEnabled {
		// signal handling has already been set up
		return
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)

	if timeout < 1 {
		Exit(Unknown, fmt.Sprintf("Invalid timeout: %d", timeout))
	}

	timedOut := time.After(time.Duration(timeout) * time.Second)
	timeoutEnabled = true

	select {
	case s := <-signals:
		Exit(Unknown, fmt.Sprintf("Received signal: %s", s))
	case <-timedOut:
		Exit(Unknown, "Timeout reached")
	}
}
