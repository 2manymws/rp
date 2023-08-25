# rp

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
