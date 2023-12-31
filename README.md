# rp [![Go Reference](https://pkg.go.dev/badge/github.com/2manymws/rp.svg)](https://pkg.go.dev/github.com/2manymws/rp) [![build](https://github.com/2manymws/rp/actions/workflows/ci.yml/badge.svg)](https://github.com/2manymws/rp/actions/workflows/ci.yml) ![Coverage](https://raw.githubusercontent.com/k1LoW/octocovs/main/badges/2manymws/rp/coverage.svg) ![Code to Test Ratio](https://raw.githubusercontent.com/k1LoW/octocovs/main/badges/2manymws/rp/ratio.svg) ![Test Execution Time](https://raw.githubusercontent.com/k1LoW/octocovs/main/badges/2manymws/rp/time.svg)

`rp` is a **r**everse **p**roxy package for multiple domains and multiple upstreams.

## Usage

Prepare an instance that implements [`rp.Relayer`](https://pkg.go.dev/github.com/2manymws/rp#Relayer) interface.

And then, create a new `http.Server` using [`rp.NewServer`](https://pkg.go.dev/github.com/2manymws/rp#NewServer) or [`rp.NewTLSServer`](https://pkg.go.dev/github.com/2manymws/rp#NewTLSServer) with the instance.

Use [`rp.NewServer`](https://pkg.go.dev/github.com/2manymws/rp#NewServer) ( [`rp.ListenAndServe`](https://pkg.go.dev/github.com/2manymws/rp#ListenAndServe) ) if handling per-domain (or per-request, as the case may be) upstreams.

```go
package main

import (
    "log"
    "net/http"

    "github.com/2manymws/rp"
)

func main() {
    var r rp.Relayer = newMyRelayer()
    log.Fatal(rp.ListenAndServe(":80", r))
}
```

Use [`rp.NewTLSServer`](https://pkg.go.dev/github.com/2manymws/rp#NewTLSServer) ( [`rp.ListenAndServeTLS`](https://pkg.go.dev/github.com/2manymws/rp#ListenAndServeTLS) )if handling per-domain TLS termination as well as per-domain HTTP request routing.

```go
package main

import (
    "log"
    "net/http"

    "github.com/2manymws/rp"
)

func main() {
    var r rp.Relayer = newMyRelayer()
    log.Fatal(rp.ListenAndServeTLS(":443", r))
}
```
