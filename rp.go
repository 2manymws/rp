package rp

import (
	"bytes"
	"crypto/tls"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
)

const errorKey = "X-Proxy-Error"

// Relayer is the interface of the implementation that determines the behavior of the reverse proxy
type Relayer interface {
	// GetUpstream returns the upstream URL for the given request.
	// DO NOT modify the request in this method.
	GetUpstream(*http.Request) (*url.URL, error)
	// Rewrite rewrites the request before sending it to the upstream.
	Rewrite(*httputil.ProxyRequest) error
	// GetCertificate returns the TLS certificate for the given client hello info.
	GetCertificate(*tls.ClientHelloInfo) (*tls.Certificate, error)
	// RoundTrip performs the round trip of the request.
	RoundTrip(r *http.Request) (*http.Response, error)
}

// NewRouter returns a new reverse proxy router.
func NewRouter(r Relayer) http.Handler {
	return &httputil.ReverseProxy{
		Rewrite: func(pr *httputil.ProxyRequest) {
			u, err := r.GetUpstream(pr.In)
			if err != nil {
				pr.Out.Header.Set(errorKey, err.Error())
				return
			}
			pr.SetURL(u)
			if err := r.Rewrite(pr); err != nil {
				pr.Out.Header.Set(errorKey, err.Error())
				return
			}
		},
		Transport: newTransport(r),
	}
}

// NewServer returns a new reverse proxy server.
func NewServer(addr string, r Relayer) *http.Server {
	rp := NewRouter(r)
	return &http.Server{
		Addr:    addr,
		Handler: rp,
	}
}

// NewTLSServer returns a new reverse proxy TLS server.
func NewTLSServer(addr string, r Relayer) *http.Server {
	rp := NewRouter(r)
	tc := new(tls.Config)
	tc.GetCertificate = r.GetCertificate
	return &http.Server{
		Addr:      addr,
		Handler:   rp,
		TLSConfig: tc,
	}
}

// ListenAndServe listens on the TCP network address addr and then proxies requests using Relayer r.
func ListenAndServe(addr string, r Relayer) error {
	s := NewServer(addr, r)
	return s.ListenAndServe()
}

// ListenAndServeTLS acts identically to ListenAndServe, except that it expects HTTPS connections.
func ListenAndServeTLS(addr string, r Relayer) error {
	s := NewTLSServer(addr, r)
	return s.ListenAndServeTLS("", "")
}

type transport struct {
	r Relayer
}

func (t *transport) RoundTrip(r *http.Request) (*http.Response, error) {
	if v := r.Header.Get(errorKey); v != "" {
		// If errorKey is set, return error response.
		body := v
		resp := &http.Response{
			Status:        http.StatusText(http.StatusBadGateway),
			StatusCode:    http.StatusBadGateway,
			Proto:         r.Proto,
			ProtoMajor:    r.ProtoMajor,
			ProtoMinor:    r.ProtoMinor,
			Body:          io.NopCloser(bytes.NewBufferString(body)),
			ContentLength: int64(len(body)),
			Request:       r,
			Header:        make(http.Header, 0),
		}
		return resp, nil
	}
	return t.r.RoundTrip(r)
}

func newTransport(r Relayer) *transport {
	return &transport{r: r}
}
