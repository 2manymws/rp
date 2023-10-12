package testutil

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path"
	"strings"
)

type Relayer struct {
	h map[string]string
}

func NewRelayer(h map[string]string) *Relayer {
	return &Relayer{
		h: h,
	}
}

func (r *Relayer) GetUpstream(req *http.Request) (*url.URL, error) {
	host := req.Host
	if upstream, ok := r.h[host]; ok {
		uu, err := url.Parse(upstream)
		if err != nil {
			return nil, err
		}
		req.URL.Scheme = uu.Scheme
		req.URL.Host = uu.Host
		req.URL.Path = strings.ReplaceAll(path.Join(uu.Path, req.URL.Path), "//", "/")
		return req.URL, nil
	}
	return nil, fmt.Errorf("not found upstream: %v", host)
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

func (r *Relayer) Rewrite(*httputil.ProxyRequest) error {
	return nil
}

type SimpleRelayer struct {
	h map[string]string
}

func NewSimpleRelayer(h map[string]string) *SimpleRelayer {
	return &SimpleRelayer{
		h: h,
	}
}

func (r *SimpleRelayer) GetUpstream(req *http.Request) (*url.URL, error) {
	host := req.Host
	if upstream, ok := r.h[host]; ok {
		uu, err := url.Parse(upstream)
		if err != nil {
			return nil, err
		}
		req.URL.Scheme = uu.Scheme
		req.URL.Host = uu.Host
		req.URL.Path = strings.ReplaceAll(path.Join(uu.Path, req.URL.Path), "//", "/")
		return req.URL, nil
	}
	return nil, fmt.Errorf("not found upstream: %v", host)
}
