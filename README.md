# rp [![Go Reference](https://pkg.go.dev/badge/github.com/k1LoW/rp.svg)](https://pkg.go.dev/github.com/k1LoW/rp) ![Coverage](https://raw.githubusercontent.com/k1LoW/octocovs/main/badges/k1LoW/rp/coverage.svg) ![Code to Test Ratio](https://raw.githubusercontent.com/k1LoW/octocovs/main/badges/k1LoW/rp/ratio.svg) ![Test Execution Time](https://raw.githubusercontent.com/k1LoW/octocovs/main/badges/k1LoW/rp/time.svg)

`rp` is a reverse proxy package for multiple domains and multiple upstreams.

## Usage

Prepare an instance that satisfies [`rp.Relayer`](https://pkg.go.dev/github.com/k1LoW/rp#Relayer) interface.

And then, create a new `http.Server` using [`rp.NewServer`](https://pkg.go.dev/github.com/k1LoW/rp#NewServer) or [`rp.NewTLSServer`](https://pkg.go.dev/github.com/k1LoW/rp#NewTLSServer) with the instance.

```go
package main

import (
    "log"
    "net/http"

    "github.com/k1LoW/rp"
)

func main() {
    r := newMyRelayer()
    s := rp.NewTLSServer(r)
    if err := s.ListenAndServe(":443"); err != nil {
        log.Fatal(err)
    }
}
```
