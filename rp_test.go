package rp

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/k1LoW/rp/testutil"
)

var _ Relayer = &testutil.Relayer{}

func TestHTTPRouting(t *testing.T) {
	_, ba := testutil.NewServer(t, "a.example.com")
	_, bb := testutil.NewServer(t, "b.example.com")
	r := testutil.NewRelayer(map[string]*url.URL{
		"a.example.com": ba,
		"b.example.com": bb,
	})

	frontend := httptest.NewServer(NewRouter(r))
	defer frontend.Close()
	tests := []struct {
		url            string
		want           string
		wantStatusCode int
	}{
		{"http://a.example.com/hello", "/hello of a.example.com", http.StatusOK},
		{"http://b.example.com/hello", "/hello of b.example.com", http.StatusOK},
		{"http://x.example.com/hello", "not found upstream: x.example.com", http.StatusBadGateway},
	}
	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			req, err := http.NewRequest("GET", tt.url, nil)
			if err != nil {
				t.Fatal(err)
			}
			u, _ := url.Parse(frontend.URL)
			req.URL.Host = u.Host
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()
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
	_, ba := testutil.NewServer(t, "a.example.com")
	_, bb := testutil.NewServer(t, "b.example.com")
	r := testutil.NewRelayer(map[string]*url.URL{
		"a.example.com": ba,
		"b.example.com": bb,
	})
	tc := new(tls.Config)
	tc.GetCertificate = r.GetCertificate
	frontend := httptest.NewUnstartedServer(NewRouter(r))
	frontend.TLS = tc
	frontend.StartTLS()
	frontend.TLS.Certificates = nil // Clear Certificates
	defer frontend.Close()
	tests := []struct {
		url            string
		want           string
		wantErr        bool
		wantStatusCode int
	}{
		{"https://a.example.com/hello", "/hello of a.example.com", false, http.StatusOK},
		{"https://b.example.com/hello", "/hello of b.example.com", false, http.StatusOK},
		{"https://x.example.com/hello", "not found upstream: x.example.com", true, 0},
	}
	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			req, err := http.NewRequest("GET", tt.url, nil)
			if err != nil {
				t.Fatal(err)
			}
			u, _ := url.Parse(frontend.URL)
			req.URL.Host = u.Host
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
						ServerName: req.Host, // Use SNI
						RootCAs:    certpool,
					},
				},
			}
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
