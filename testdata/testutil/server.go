package testutil

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func NewUpstreamServer(t *testing.T, host string) *httptest.Server {
	t.Helper()
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(fmt.Sprintf("response from %s [%s]", r.URL.String(), host)))
	})
	ts := httptest.NewServer(h)
	t.Cleanup(ts.Close)
	return ts
}
