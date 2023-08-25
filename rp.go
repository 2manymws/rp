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

type Relayer interface {
	GetUpstream(*http.Request) (*url.URL, error)
	Rewrite(*httputil.ProxyRequest) error
	GetCertificate(*tls.ClientHelloInfo) (*tls.Certificate, error)
	RoundTrip(r *http.Request) (*http.Response, error)
}

func NewRouter(r Relayer) http.Handler {
	return &httputil.ReverseProxy{
		Rewrite: func(pr *httputil.ProxyRequest) {
			u, err := r.GetUpstream(pr.In)
			if err != nil {
				pr.Out.Header.Set(errorKey, err.Error())
				return
			}
			pr.SetURL(u)
			pr.SetXForwarded()
			if err := r.Rewrite(pr); err != nil {
				pr.Out.Header.Set(errorKey, err.Error())
				return
			}
		},
		Transport: newTransport(r),
	}
}

func NewServer(addr string, r Relayer) *http.Server {
	rp := NewRouter(r)
	return &http.Server{
		Addr:    addr,
		Handler: rp,
	}
}

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
