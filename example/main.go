package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/k1LoW/rp"
)

type myRelayer struct {
}

func (r *myRelayer) GetUpstream(req *http.Request) (*url.URL, error) {
	u, _ := url.Parse("https://httpbin.org")
	u.Path = "/anything" + req.URL.Path
	return u, nil
}
func (r *myRelayer) Rewrite(pr *httputil.ProxyRequest) error {
	pr.Out.Host = "httpbin.org"
	return nil
}

func (r *myRelayer) GetCertificate(*tls.ClientHelloInfo) (*tls.Certificate, error) {
	return nil, nil
}

func (r *myRelayer) RoundTrip(req *http.Request) (*http.Response, error) {
	return http.DefaultTransport.RoundTrip(req)
}

func main() {
	var r rp.Relayer = &myRelayer{}
	log.Fatal(rp.ListenAndServe(":8088", r))
}
