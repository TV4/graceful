/*

Package graceful simplifies graceful shutdown of HTTP servers (Go 1.8+)

Installation

Just go get the package:

    go get -u github.com/TV4/graceful

Usage

A small usage example

    package main

    import (
    	"log"
    	"net/http"
    	"time"

    	"github.com/TV4/graceful"
    )

    type server struct{}

    func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    	time.Sleep(2 * time.Second)
    	w.Write([]byte("Hello!"))
    }

    func main() {
    	addr := ":2017"

    	log.Printf("Listening on http://0.0.0.0%s\n", addr)

    	graceful.ListenAndServe(&http.Server{
    		Addr:    addr,
    		Handler: &server{},
    	})
    }

*/
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

// Server is implemented by *http.Server
type Server interface {
	ListenAndServe() error
	Shutdowner
}

// Shutdowner is implemented by *http.Server
type Shutdowner interface {
	Shutdown(ctx context.Context) error
}

// Logger is implemented by *log.Logger
type Logger interface {
	Printf(format string, v ...interface{})
	Fatal(...interface{})
}

// DefaultTimeout for context used in call to *http.Server.Shutdown
var DefaultTimeout = 15 * time.Second

// DefaultLogger is the logger used by the shutdown function
var DefaultLogger Logger = log.New(os.Stdout, "", 0)

// Format strings used by the logger
var (
	ListeningFormat = "Listening on http://0.0.0.0%s\n"
	ShutdownFormat  = "\nShutdown with timeout: %s\n"
	ErrorFormat     = "Error: %v\n"
	StoppedFormat   = "Server stopped\n"
)

// LogListenAndServe logs using the DefaultLogger and then calls ListenAndServe
func LogListenAndServe(hs *http.Server) {
	DefaultLogger.Printf(ListeningFormat, hs.Addr)

	ListenAndServe(hs)
}

// ListenAndServe starts the server in a goroutine and then calls Shutdown
func ListenAndServe(s Server) {
	go func() {
		if err := s.ListenAndServe(); err != http.ErrServerClosed {
			DefaultLogger.Fatal(err)
		}
	}()

	Shutdown(s)
}

// Shutdown blocks until os.Interrupt or syscall.SIGTERM received, then
// running *http.Server.Shutdown with a context having a timeout
func Shutdown(s Shutdowner) {
	wait()

	shutdown(s, DefaultLogger, DefaultTimeout)
}

func wait() {
	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
}

func shutdown(s Shutdowner, logger Logger, timeout time.Duration) {
	if s == nil {
		return
	}

	if logger == nil {
		logger = log.New(ioutil.Discard, "", 0)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	logger.Printf(ShutdownFormat, timeout)

	if err := s.Shutdown(ctx); err != nil {
		logger.Printf(ErrorFormat, err)
	} else {
		logger.Printf(StoppedFormat)
	}
}
