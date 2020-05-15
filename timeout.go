package check

import (
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Start the timeout and signal handler in a goroutine
func (f *Flags) EnableTimeoutHandler() {
	go HandleTimeout(f.Timeout)
}

// Helper for a goroutine, to wait for signals and timeout, and exit with a proper code
func HandleTimeout(timeout int) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)

	if timeout < 1 {
		Exit(Unknown, "Invalid timeout: %d", timeout)
	}

	timedOut := time.After(time.Duration(timeout) * time.Second)

	select {
	case s := <-signals:
		Exit(Unknown, "Received signal: %s", s)
	case <-timedOut:
		Exit(Unknown, "Timeout reached")
	}
}
