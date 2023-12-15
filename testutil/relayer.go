package testutil

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
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
		ru, err := url.Parse(req.URL.String())
		if err != nil {
			return nil, err
		}
		uu, err := url.Parse(upstream)
		if err != nil {
			return nil, err
		}
		ru.Scheme = uu.Scheme
		ru.Host = uu.Host
		ru.Path = strings.ReplaceAll(path.Join(uu.Path, req.URL.Path), "//", "/")
		return ru, nil
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

type RoundTripOnErrorRelayer struct {
	SimpleRelayer
}

func NewRoundTripOnErrorRelayer(h map[string]string) *RoundTripOnErrorRelayer {
	return &RoundTripOnErrorRelayer{
		SimpleRelayer: SimpleRelayer{
			h: h,
		},
	}
}
func (r *RoundTripOnErrorRelayer) RoundTripOnError(req *http.Request) (*http.Response, error) {
	body := fmt.Sprintf("round trip on error: %v", req.Host)
	return &http.Response{
		Status:        http.StatusText(http.StatusBadGateway),
		StatusCode:    http.StatusBadGateway,
		Proto:         req.Proto,
		ProtoMajor:    req.ProtoMajor,
		ProtoMinor:    req.ProtoMinor,
		Body:          io.NopCloser(bytes.NewBufferString(body)),
		ContentLength: int64(len(body)),
		Request:       req,
		Header:        make(http.Header, 0),
	}, nil
}
