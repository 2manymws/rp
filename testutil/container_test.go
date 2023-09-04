package testutil

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestContainer(t *testing.T) {
	_ = NewUpstreamEchoNGINXServer(t, "a.example.com")
	upstreams := map[string]string{
		"a.example.com": fmt.Sprintf("http://%s:80", "a.example.com"),
	}
	proxy := NewReverseProxyNGINXServer(t, "r.example.com", upstreams)
	{
		now := time.Now()
		req, err := http.NewRequest("GET", proxy+"sleep", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Host = "a.example.com"
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		after := time.Now()
		if after.Sub(now) < 1*time.Second {
			t.Fatal("sleep.js is not working")
		}
	}

	{
		req, err := http.NewRequest("GET", proxy+"sleep", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Host = "a.example.com"
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		if res.Header.Get("X-Nginx-Cache") != "HIT" {
			t.Fatal("NGINX cache is not working")
		}
	}
}
