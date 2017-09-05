package graceful

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"testing"
)

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

		want := fmt.Sprintf(ShutdownFormat+ServerStopped, Timeout)

		if got := buf.String(); got != want {
			t.Fatalf("buf.String() = %q, want %q", got, want)
		}
	})
}
