package rp_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/k1LoW/rp"
	"github.com/k1LoW/rp/testutil"
)

func BenchmarkNGINX(b *testing.B) {
	var upstreams = map[string]string{
		"a.example.com": "",
		"b.example.com": "",
		"c.example.com": "",
	}
	for hostname := range upstreams {
		_ = testutil.NewUpstreamEchoNGINXServer(b, hostname)
		upstreams[hostname] = fmt.Sprintf("http://%s:80", hostname)
	}
	proxy := testutil.NewReverseProxyNGINXServer(b, "r.example.com", upstreams)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			upstream := sample(upstreams)
			req, err := http.NewRequest("GET", proxy, nil)
			if err != nil {
				b.Error(err)
				return
			}
			req.Host = upstream
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				b.Error(err)
				return
			}
			got := res.StatusCode
			want := http.StatusOK
			if res.StatusCode != http.StatusOK {
				b.Errorf("got %v want %v", got, want)
			}
		}
	})
}

func BenchmarkRP(b *testing.B) {
	var upstreams = map[string]string{
		"a.example.com": "",
		"b.example.com": "",
		"c.example.com": "",
	}
	for hostname := range upstreams {
		urlstr := testutil.NewUpstreamEchoNGINXServer(b, hostname)
		upstreams[hostname] = urlstr
	}
	r := testutil.NewRelayer(upstreams)
	proxy := httptest.NewServer(rp.NewRouter(r))
	defer proxy.Close()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			upstream := sample(upstreams)
			req, err := http.NewRequest("GET", proxy.URL, nil)
			if err != nil {
				b.Error(err)
				return
			}
			req.Host = upstream
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				b.Error(err)
				return
			}
			got := res.StatusCode
			want := http.StatusOK
			if res.StatusCode != http.StatusOK {
				b.Errorf("got %v want %v", got, want)
			}
		}
	})
}

func sample[T any](m map[string]T) string {
	for k := range m {
		return k
	}
	return ""
}
