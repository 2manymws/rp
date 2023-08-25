package testutil

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func NewServer(t *testing.T, h string) (*httptest.Server, *url.URL) {
	t.Helper()
	u, err := url.Parse(fmt.Sprintf("http://%s", h))
	if err != nil {
		t.Fatal(err)
	}
	u.Scheme = "http"
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(fmt.Sprintf("/ of %s", h)))
	})
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(fmt.Sprintf("/hello of %s", h)))
	})
	s := httptest.NewServer(mux)
	t.Cleanup(s.Close)
	u.Host = s.Listener.Addr().String()
	return s, u
}
