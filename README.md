# HTTP ROUTER

A very naive implementation of http router that supports path param parsing that is suitable
for most REST implementations


## Usage

```go
package main

import (
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
    		w.Write([]byte("Pong!"))
    	})
	
    log.Fatal(http.ListenAndServe(":8080", rtr))
}
```

## Running tests

```bash
> go test -v ./...

```


## Running Benchmarks

```bash
> go test -run none -bench Benchmark -benchmem -benchtime 3s -memprofile mem.out -c
```

