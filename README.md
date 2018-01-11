# graceful

[![Build Status](https://travis-ci.org/TV4/graceful.svg?branch=master)](https://travis-ci.org/TV4/graceful)
[![Go Report Card](https://goreportcard.com/badge/github.com/TV4/graceful)](https://goreportcard.com/report/github.com/TV4/graceful)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/TV4/graceful)
[![License MIT](https://img.shields.io/badge/license-MIT-lightgrey.svg?style=flat)](https://github.com/TV4/graceful#license-mit)

## Installation

    go get -u github.com/TV4/graceful

## Usage

```go
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
```

```
$ go run main.go
2017/06/19 16:35:28 Listening on http://0.0.0.0:2017
^C
```

### You can also use `graceful.LogListenAndServe`

```go
package main

import (
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
	graceful.LogListenAndServe(&http.Server{
		Addr:    ":2017",
		Handler: &server{},
	})
}
```

```
$ go run main.go
Listening on http://0.0.0.0:2017
^C
Server shutdown with timeout: 15s
Finished all in-flight HTTP requests
Shutdown finished 14s before deadline
```

### And optionally your handler can implement the Shutdowner interface

```go
package main

import (
	"context"
	"fmt"
	"net/http"

	graceful "github.com/TV4/graceful"
)

func main() {
	graceful.LogListenAndServe(&http.Server{
		Addr:    ":8080",
		Handler: &server{},
	})
}

type server struct{}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world!")
}

func (s *server) Shutdown(ctx context.Context) error {
	// Here you can do anything that you need to do
	// after the *http.Server has stopped accepting
	// new connections and finished its Shutdown

	// The ctx is the same as the context used to
	// perform *https.Server.Shutdown and thus
	// shares the timeout (15 seconds by default)

	fmt.Println("Finished *server.Shutdown")

	return nil
}
```

```
$ go run main.go
Listening on http://0.0.0.0:8080
^C
Server shutdown with timeout: 15s
Finished all in-flight HTTP requests
Shutting down handler with timeout: 15s
Finished *server.Shutdown
Shutdown finished 15s before deadline
```

## License (MIT)

Copyright (c) 2017-2018 TV4

> Permission is hereby granted, free of charge, to any person obtaining
> a copy of this software and associated documentation files (the
> "Software"), to deal in the Software without restriction, including
> without limitation the rights to use, copy, modify, merge, publish,
> distribute, sublicense, and/or sell copies of the Software, and to
> permit persons to whom the Software is furnished to do so, subject to
> the following conditions:

> The above copyright notice and this permission notice shall be
> included in all copies or substantial portions of the Software.

> THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
> EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
> MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
> NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
> LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
> OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
> WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
