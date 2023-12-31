package main

import (
	"log"
	"net/url"

	"github.com/2manymws/rp"
	"github.com/2manymws/rp/testutil"
)

func main() {
	r := testutil.NewRelayer(
		map[string]*url.URL{
			"127.0.0.1:18082": &url.URL{
				Scheme: "http",
				Host:   "127.0.0.1:18080",
			},
		},
	)
	s := rp.NewServer("127.0.0.1:18082", r)
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}
