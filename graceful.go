package graceful

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Timeout for context used in call to *http.Server.Shutdown
var Timeout = 15 * time.Second

// Format strings used by the logger
var (
	ShutdownFormat = "\nShutdown with timeout: %s\n"
	ErrorFormat    = "Error: %v\n"
	StoppedFormat  = "Server stopped\n"
)

// Shutdown blocks until os.Interrupt or syscall.SIGTERM received, then
// running *http.Server.Shutdown with a context having a timeout
func Shutdown(hs *http.Server, logger *log.Logger) {
	wait()

	shutdown(hs, logger, Timeout)
}

func wait() {
	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
}

func shutdown(hs *http.Server, logger *log.Logger, timeout time.Duration) {
	if hs == nil {
		return
	}

	if logger == nil {
		logger = log.New(ioutil.Discard, "", 0)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	logger.Printf(ShutdownFormat, timeout)

	if err := hs.Shutdown(ctx); err != nil {
		logger.Printf(ErrorFormat, err)
	} else {
		logger.Printf(StoppedFormat)
	}
}
