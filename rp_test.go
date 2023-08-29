package rp_test

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/k1LoW/rp"
	"github.com/k1LoW/rp/testutil"
)

var _ rp.Relayer = &testutil.Relayer{}

type upstream struct {
	hostname string
	rootPath string
}

func TestHTTPRouting(t *testing.T) {
	tests := []struct {
		upstreams      []upstream
		reqURL         string
		want           string
		wantStatusCode int
	}{
		{[]upstream{{"a.example.com", "/"}, {"b.example.com", "/"}}, "http://a.example.com", "response from / [a.example.com]", http.StatusOK},
		{[]upstream{{"a.example.com", "/"}, {"b.example.com", "/"}}, "http://b.example.com/hello", "response from /hello [b.example.com]", http.StatusOK},
		{[]upstream{{"a.example.com", "/"}, {"b.example.com", "/"}}, "http://a.example.com/hello?foo=bar", "response from /hello?foo=bar [a.example.com]", http.StatusOK},
		{[]upstream{{"a.example.com", "/"}, {"b.example.com", "/"}}, "http://x.example.com/hello", "not found upstream: x.example.com", http.StatusBadGateway},
		{[]upstream{{"a.example.com", "/A"}}, "http://a.example.com", "response from /A [a.example.com]", http.StatusOK},
		{[]upstream{{"a.example.com", "/A"}}, "http://a.example.com/hello?foo=bar", "response from /A/hello?foo=bar [a.example.com]", http.StatusOK},
	}
	for _, tt := range tests {
		t.Run(tt.reqURL, func(t *testing.T) {
			// Arrange
			m := map[string]string{}
			for _, u := range tt.upstreams {
				us := testutil.NewUpstreamServer(t, u.hostname)
				m[u.hostname] = us.URL + u.rootPath
			}
			r := testutil.NewRelayer(m)
			port, err := testutil.NewPort()
			if err != nil {
				t.Fatal(err)
			}
			proxy := rp.NewServer(fmt.Sprintf("127.0.0.1:%d", port), r)
			go func() {
				t.Helper()
				if err := proxy.ListenAndServe(); err != nil {
					if !errors.Is(err, http.ErrServerClosed) {
						t.Error(err)
					}
				}
			}()
			t.Cleanup(func() {
				_ = proxy.Shutdown(context.Background())
			})
			t.Cleanup(func() {
				_ = proxy.Shutdown(context.Background())
			})
			proxyURL, err := url.Parse(fmt.Sprintf("http://127.0.0.1:%d", port))
			if err != nil {
				t.Fatal(err)
			}
			client := http.DefaultClient
			for {
				if _, err := client.Get(proxyURL.String()); err == nil {
					break
				}
				time.Sleep(100 * time.Millisecond)
			}

			// Act
			req, err := http.NewRequest("GET", tt.reqURL, nil)
			if err != nil {
				t.Fatal(err)
			}
			req.URL.Host = proxyURL.Host
			resp, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			// Assert
			b, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}
			got := string(b)
			if got != tt.want {
				t.Errorf("got %v\nwant %v", got, tt.want)
			}
			if resp.StatusCode != tt.wantStatusCode {
				t.Errorf("got %v\nwant %v", resp.StatusCode, tt.wantStatusCode)
			}
		})
	}
}

func TestHTTPSRouting(t *testing.T) {
	tests := []struct {
		upstreams      []upstream
		reqURL         string
		want           string
		wantErr        bool
		wantStatusCode int
	}{
		{[]upstream{{"a.example.com", "/"}, {"b.example.com", "/"}}, "https://a.example.com/hello", "response from /hello [a.example.com]", false, http.StatusOK},
		{[]upstream{{"a.example.com", "/"}, {"b.example.com", "/"}}, "https://b.example.com/hello", "response from /hello [b.example.com]", false, http.StatusOK},
		{[]upstream{{"a.example.com", "/"}, {"b.example.com", "/"}}, "https://x.example.com/hello", "not found upstream: x.example.com", true, 0},
	}
	for _, tt := range tests {
		t.Run(tt.reqURL, func(t *testing.T) {
			// Arrange
			m := map[string]string{}
			for _, u := range tt.upstreams {
				us := testutil.NewUpstreamServer(t, u.hostname)
				m[u.hostname] = us.URL + u.rootPath
			}
			r := testutil.NewRelayer(m)
			port, err := testutil.NewPort()
			if err != nil {
				t.Fatal(err)
			}
			proxy := rp.NewTLSServer(fmt.Sprintf("127.0.0.1:%d", port), r)
			go func() {
				t.Helper()
				if err := proxy.ListenAndServeTLS("", ""); err != nil {
					if !errors.Is(err, http.ErrServerClosed) {
						t.Error(err)
					}
				}
			}()
			t.Cleanup(func() {
				_ = proxy.Shutdown(context.Background())
			})
			proxyURL, err := url.Parse(fmt.Sprintf("https://127.0.0.1:%d", port))
			if err != nil {
				t.Fatal(err)
			}
			certpool, err := x509.SystemCertPool()
			if err != nil {
				// FIXME for Windows
				// ref: https://github.com/golang/go/issues/18609
				certpool = x509.NewCertPool()
			}
			cacert, err := os.ReadFile("testdata/cacert.pem")
			if err != nil {
				t.Fatal(err)
			}
			if !certpool.AppendCertsFromPEM(cacert) {
				t.Fatal("failed to add cacert")
			}
			client := http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{
						ServerName: "a.example.com",
						RootCAs:    certpool,
					},
					ForceAttemptHTTP2: true,
				},
			}
			for {
				if _, err := client.Get(proxyURL.String()); err == nil {
					break
				}
				time.Sleep(100 * time.Millisecond)
			}

			// Act
			req, err := http.NewRequest("GET", tt.reqURL, nil)
			if err != nil {
				t.Fatal(err)
			}
			req.URL.Host = proxyURL.Host
			client.CloseIdleConnections()
			client.Transport.(*http.Transport).TLSClientConfig.ServerName = req.Host // Use SNI
			resp, err := client.Do(req)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("got error: %v", err)
				}
				return
			}
			if tt.wantErr {
				t.Error("want error")
				return
			}
			defer resp.Body.Close()

			// Assert
			b, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}
			if got := string(b); got != tt.want {
				t.Errorf("got %v\nwant %v", got, tt.want)
			}
			if got := resp.Proto; got != "HTTP/2.0" {
				t.Errorf("got %v\nwant %v", got, "HTTP/2.0")
			}
			if got := resp.StatusCode; got != tt.wantStatusCode {
				t.Errorf("got %v\nwant %v", resp.StatusCode, tt.wantStatusCode)
			}
		})
	}
}
