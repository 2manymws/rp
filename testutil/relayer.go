package testutil

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

type Relayer struct {
	h map[string]*url.URL
}

func NewRelayer(h map[string]*url.URL) *Relayer {
	return &Relayer{
		h: h,
	}
}

func (r *Relayer) GetUpstream(req *http.Request) (*url.URL, error) {
	host := req.Host
	if dest, ok := r.h[host]; ok {
		return dest, nil
	}
	return nil, fmt.Errorf("not found upstream: %v", host)
}

func (r *Relayer) Rewrite(pr *httputil.ProxyRequest) error {
	pr.SetXForwarded()
	return nil
}

func (r *Relayer) GetCertificate(i *tls.ClientHelloInfo) (*tls.Certificate, error) {
	cert := fmt.Sprintf("testdata/%s.cert.pem", i.ServerName)
	key := fmt.Sprintf("testdata/%s.key.pem", i.ServerName)
	if _, err := os.Stat(cert); err != nil {
		return nil, err
	}
	if _, err := os.Stat(key); err != nil {
		return nil, err
	}
	c, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *Relayer) RoundTrip(req *http.Request) (*http.Response, error) {
	return http.DefaultTransport.RoundTrip(req)
}
