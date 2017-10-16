package graceful

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

func TestListenAndServeTLS(t *testing.T) {
	var buf bytes.Buffer

	logger := log.New(&buf, "", 0)

	go func() {
		time.Sleep(10 * time.Millisecond)
		signals <- os.Interrupt
	}()

	ListenAndServeTLS(&http.Server{
		Addr: ":0", Handler: &testHandler{logger},
	}, "testdata/server.crt", "testdata/server.key")

	s := buf.String()

	for _, want := range []string{
		"Shutdown in testHandler",
	} {
		if !strings.Contains(s, want) {
			t.Fatalf("log output does not include %q", want)
		}
	}
}

func TestLogListenAndServe(t *testing.T) {
	t.Run("with logger", func(t *testing.T) {
		var buf bytes.Buffer

		logger := log.New(&buf, "", 0)

		go func() {
			time.Sleep(10 * time.Millisecond)
			signals <- os.Interrupt
		}()

		LogListenAndServe(&http.Server{
			Addr: ":0", Handler: &testHandler{logger},
		}, logger)

		s := buf.String()

		for _, want := range []string{
			"Server shutdown with timeout: 15s",
			"Shutdown in testHandler",
		} {
			if !strings.Contains(s, want) {
				t.Fatalf("log output does not include %q", want)
			}
		}
	})

	t.Run("with no logger", func(t *testing.T) {
		go func() {
			time.Sleep(10 * time.Millisecond)
			signals <- os.Interrupt
		}()

		LogListenAndServe(&http.Server{
			Addr: ":0", Handler: &testHandler{},
		})
	})

	t.Run("with nil logger", func(t *testing.T) {
		go func() {
			time.Sleep(10 * time.Millisecond)
			signals <- os.Interrupt
		}()

		LogListenAndServe(&http.Server{
			Addr: ":0", Handler: &testHandler{},
		}, nil)
	})
}

func TestShutdown(t *testing.T) {
	t.Run("nil-hs", func(t *testing.T) {
		shutdown(nil, nil)
	})

	t.Run("nil-logger", func(t *testing.T) {
		shutdown(&http.Server{}, nil)
	})

	t.Run("logger", func(t *testing.T) {
		var buf bytes.Buffer

		shutdown(&http.Server{}, log.New(&buf, "", 0))

		want := fmt.Sprintf(ShutdownFormat+FinishedHTTP+FinishedFormat, Timeout, 15)

		if got := buf.String(); got != want {
			t.Fatalf("buf.String() = %q, want %q", got, want)
		}
	})
}

type testHandler struct {
	logger *log.Logger
}

func (h *testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {}

func (h *testHandler) Shutdown(ctx context.Context) error {
	if h.logger != nil {
		h.logger.Println("Shutdown in testHandler")
	}

	return nil
}
