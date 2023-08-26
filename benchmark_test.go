package rp

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/k1LoW/rp/testutil"
)

func BenchmarkNGINX(b *testing.B) {
	var upstreams = map[string]string{
		"a.example.com": "",
		"b.example.com": "",
		"c.example.com": "",
	}
	for hostname := range upstreams {
		_ = testutil.CreateUpstreamEchoServer(b, hostname)
		upstreams[hostname] = fmt.Sprintf("http://%s:80", hostname)
	}
	proxy := testutil.CreateReverseProxyServer(b, "r.example.com", upstreams)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		upstream := sample(upstreams)
		req, err := http.NewRequest("GET", proxy, nil)
		if err != nil {
			b.Fatal(err)
		}
		req.Host = upstream
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			b.Fatal(err)
		}
		if res.StatusCode != http.StatusOK {
			b.Fatal(res.Status)
		}
	}
}

func BenchmarkRP(b *testing.B) {
	var upstreams = map[string]*url.URL{
		"a.example.com": nil,
		"b.example.com": nil,
		"c.example.com": nil,
	}
	for hostname := range upstreams {
		host := testutil.CreateUpstreamEchoServer(b, hostname)
		u, err := url.Parse(host)
		if err != nil {
			b.Fatal(err)
		}
		upstreams[hostname] = u
	}
	r := testutil.NewRelayer(upstreams)
	proxy := httptest.NewServer(NewRouter(r))
	defer proxy.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		upstream := sample(upstreams)
		req, err := http.NewRequest("GET", proxy.URL, nil)
		if err != nil {
			b.Fatal(err)
		}
		req.Host = upstream
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			b.Fatal(err)
		}
		if res.StatusCode != http.StatusOK {
			b.Fatal(res.Status)
		}
	}
}

func sample[T any](m map[string]T) string {
	for k := range m {
		return k
	}
	return ""
}
