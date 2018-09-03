# HTTP ROUTER

A very naive implementation of http router that supports path param parsing, suitable
for most REST implementations. Inspired by [httprouter](https://github.com/julienschmidt/httprouter/blob/master/params_go17.go)

## Usage

```go
package main

import (
	"fmt"
	"github.com/shyamz-22/router"
	"log"
	"net/http"
)

func main() {
	
	rtr := router.New()
	
	rtr.Add("/ping", http.MethodGet, func(w http.ResponseWriter, r *http.Request, params router.PathParams) {
		w.Write([]byte("Pong!"))
	})
	
	rtr.Add("/pings/:id", http.MethodGet, func(w http.ResponseWriter, r *http.Request, params router.PathParams) {
    	id := params.ByName("id")	
    	response := fmt.Sprintf("Pong: %s", id)
		w.Write([]byte(response))
    	})
	
    log.Fatal(http.ListenAndServe(":8080", rtr))
}
```

## Running tests

```bash
> go test -v ./...

```

## Running parallel tests

```bash
> go test -parallel 8 -v ./...

```


## Running Benchmarks

```bash
> go test -run none -bench Benchmark -benchmem -benchtime 3s -memprofile mem.out
```

## Memory profiling

```bash
> go test -run none -bench BenchmarkGithub -benchmem -benchtime 20s -memprofile mem.out
> go tool pprof -alloc_space router.test mem.out
 
```
### Output

```
File: router.test
Type: alloc_space
Time: Sep 2, 2018 at 9:32pm (CEST)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
(pprof) list findRoute

```

## Setup

```bash
> go get github.com/shyamz-22/router
> cd $GOPATH/github.com/shyamz-22/router
> dep ensure -update
> go test ./...
```
