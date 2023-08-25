# rp ![Coverage](https://raw.githubusercontent.com/k1LoW/octocovs/main/badges/k1LoW/rp/coverage.svg) ![Code to Test Ratio](https://raw.githubusercontent.com/k1LoW/octocovs/main/badges/k1LoW/rp/ratio.svg) ![Test Execution Time](https://raw.githubusercontent.com/k1LoW/octocovs/main/badges/k1LoW/rp/time.svg)

`rp` is a reverse proxy package for multiple domains and multiple upstreams.

## Usage

First, prepare an instance that satisfies `rp.Relayer` interface.

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
