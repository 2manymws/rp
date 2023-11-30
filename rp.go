package rp

import (
	"bytes"
	"crypto/tls"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

const errorKey = "X-Proxy-Error"

// Relayer is the interface of the implementation that determines the behavior of the reverse proxy
type Relayer interface { //nostyle:ifacenames
	// GetUpstream returns the upstream URL for the given request.
	// If upstream is not determined, nil may be returned
	// DO NOT modify the request in this method.
	GetUpstream(*http.Request) (*url.URL, error)
}

type Rewriter interface {
	// Rewrite rewrites the request before sending it to the upstream.
	// For example, you can set `X-Forwarded-*` header here using [httputil.ProxyRequest.SetXForwarded](https://pkg.go.dev/net/http/httputil#ProxyRequest.SetXForwarded)
	Rewrite(*httputil.ProxyRequest) error
}

type CertGetter interface { //nostyle:ifacenames
	// GetCertificate returns the TLS certificate for the given client hello info.
	GetCertificate(*tls.ClientHelloInfo) (*tls.Certificate, error)
}

type RoundTripper interface {
	// RoundTrip performs the round trip of the request.
	// It is necessary to implement the functions that http.Transport is responsible for (e.g. MaxIdleConnsPerHost).
	RoundTrip(r *http.Request) (*http.Response, error)
	RoundTripOnError(r *http.Request) (*http.Response, error)
}

type RoundTipperOnErrorer interface {
	// RoundTripOnError performs the round trip of the request when the upstream returns an error.
	// If this method is not implemented, the request will be sent to the default transport error.
	RoundTripOnError(r *http.Request) (*http.Response, error)
}

type ErrorHandler interface {
	// ErrorHandler handles the error returned by the upstream.
	// If this method is not implemented, the error will be returned to the client.
	ErrorHandler(http.ResponseWriter, *http.Request, error)
}
type relayer struct {
	Relayer
	Rewrite          func(*httputil.ProxyRequest) error
	GetCertificate   func(*tls.ClientHelloInfo) (*tls.Certificate, error)
	RoundTrip        func(*http.Request) (*http.Response, error)
	RoundTripOnError func(*http.Request) (*http.Response, error)
	ErrorHandler     func(http.ResponseWriter, *http.Request, error)
}

func newRelayer(r Relayer) *relayer {
	rr := &relayer{
		Relayer: r,
	}
	if v, ok := r.(Rewriter); ok {
		rr.Rewrite = v.Rewrite
	}
	if v, ok := r.(CertGetter); ok {
		rr.GetCertificate = v.GetCertificate
	}
	if v, ok := r.(RoundTripper); ok {
		rr.RoundTrip = v.RoundTrip
	} else {
		rr.RoundTrip = http.DefaultTransport.RoundTrip
	}
	if v, ok := r.(RoundTipperOnErrorer); ok {
		rr.RoundTripOnError = v.RoundTripOnError
	}
	if v, ok := r.(ErrorHandler); ok {
		rr.ErrorHandler = v.ErrorHandler
	}
	return rr
}

// NewRouter returns a new reverse proxy router.
func NewRouter(r Relayer) http.Handler {
	rr := newRelayer(r)
	if rr.Rewrite == nil {
		return &httputil.ReverseProxy{
			Rewrite: func(pr *httputil.ProxyRequest) {
				u, err := rr.GetUpstream(pr.In)
				if err != nil {
					pr.Out.Header.Set(errorKey, err.Error())
					return
				}
				if u != nil {
					if strings.HasPrefix(u.Host, "/") {
						pr.Out.Host = pr.In.Host
					} else {
						pr.Out.Host = u.Host
					}
					pr.Out.URL = u
					pr.SetXForwarded()
				}
			},
			Transport:    newTransport(rr),
			ErrorHandler: rr.ErrorHandler,
		}
	}
	return &httputil.ReverseProxy{
		Rewrite: func(pr *httputil.ProxyRequest) {
			u, err := rr.GetUpstream(pr.In)
			if err != nil {
				pr.Out.Header.Set(errorKey, err.Error())
				return
			}
			if u != nil {
				if strings.HasPrefix(u.Host, "/") {
					pr.Out.Host = pr.In.Host
				} else {
					pr.Out.Host = u.Host
				}

				pr.Out.URL = u
			}
			if err := rr.Rewrite(pr); err != nil {
				pr.Out.Header.Set(errorKey, err.Error())
				return
			}
		},
		Transport:    newTransport(rr),
		ErrorHandler: rr.ErrorHandler,
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
	rr := newRelayer(r)
	tc := new(tls.Config)
	if rr.GetCertificate != nil {
		tc.GetCertificate = rr.GetCertificate
	}
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
	rr *relayer
}

func (t *transport) RoundTrip(r *http.Request) (*http.Response, error) {
	if v := r.Header.Get(errorKey); v != "" {
		if t.rr.RoundTripOnError != nil {
			return t.rr.RoundTripOnError(r)
		}
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
	return t.rr.RoundTrip(r)
}

func newTransport(rr *relayer) *transport {
	return &transport{rr: rr}
}
