# rp [![Go Reference](https://pkg.go.dev/badge/github.com/k1LoW/rp.svg)](https://pkg.go.dev/github.com/k1LoW/rp) ![Coverage](https://raw.githubusercontent.com/k1LoW/octocovs/main/badges/k1LoW/rp/coverage.svg) ![Code to Test Ratio](https://raw.githubusercontent.com/k1LoW/octocovs/main/badges/k1LoW/rp/ratio.svg) ![Test Execution Time](https://raw.githubusercontent.com/k1LoW/octocovs/main/badges/k1LoW/rp/time.svg)

`rp` is a reverse proxy package for multiple domains and multiple upstreams.

## Usage

Prepare an instance that satisfies [`rp.Relayer`](https://pkg.go.dev/github.com/k1LoW/rp#Relayer) interface.

And then, create a new `http.Server` using [`rp.NewServer`](https://pkg.go.dev/github.com/k1LoW/rp#NewServer) or [`rp.NewTLSServer`](https://pkg.go.dev/github.com/k1LoW/rp#NewTLSServer) with the instance.

```go
package main

import (
    "errors"
    "log"
    "net/http"

    "github.com/k1LoW/rp"
)

func main() {
    var r rp.Relayer = newMyRelayer()
    if err := rp.ListenAndServe(":80", r); err != nil {
        if !errors.Is(err, http.ErrServerClosed) {
            log.Fatal(err)
        }
    }
}
```

```go
package main

import (
    "errors"
    "log"
    "net/http"

    "github.com/k1LoW/rp"
)

func main() {
    var r rp.Relayer = newMyRelayer()
    if err := rp.ListenAndServeTLS(":443", r); err != nil {
        if !errors.Is(err, http.ErrServerClosed) {
            log.Fatal(err)
        }
    }
}
```
