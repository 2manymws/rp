package testutil

import (
	"embed"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"text/template"

	"github.com/ory/dockertest/v3"
)

const (
	networkName = "rp-test-network"
)

//go:embed templates/*
var conf embed.FS

func NewReverseProxyNGINXServer(t testing.TB, hostname string, upstreams map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	tb, err := conf.ReadFile("templates/nginx_reverse.conf.tmpl")
	if err != nil {
		t.Fatal(err)
	}
	tmpl := template.Must(template.New("conf").Parse(string(tb)))
	p := filepath.Join(dir, fmt.Sprintf("%s.nginx_reverse.conf", hostname))
	f, err := os.Create(p)
	if err != nil {
		t.Fatal(err)
	}
	if err := tmpl.Execute(f, map[string]any{
		"Hostname":  hostname,
		"Upstreams": upstreams,
	}); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
	return createNGINXServer(t, hostname, p)
}

func NewUpstreamEchoNGINXServer(t testing.TB, hostname string) string {
	t.Helper()
	dir := t.TempDir()
	tb, err := conf.ReadFile("templates/nginx_echo.conf.tmpl")
	if err != nil {
		t.Fatal(err)
	}
	tmpl := template.Must(template.New("conf").Parse(string(tb)))
	p := filepath.Join(dir, fmt.Sprintf("%s.nginx_echo.conf", hostname))
	f, err := os.Create(p)
	if err != nil {
		t.Fatal(err)
	}
	if err := tmpl.Execute(f, map[string]any{
		"Hostname": hostname,
	}); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
	return createNGINXServer(t, hostname, p)
}

func createNGINXServer(t testing.TB, hostname, confp string) string {
	t.Helper()
	dir := t.TempDir()
	sb, err := conf.ReadFile("templates/sleep.js")
	if err != nil {
		t.Fatal(err)
	}
	sp := filepath.Join(dir, "sleep.js")
	if err := os.WriteFile(sp, sb, 0644); err != nil {
		t.Fatal(err)
	}
	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("Could not connect to docker: %s", err)
	}
	opt := &dockertest.RunOptions{
		Hostname:   hostname,
		Repository: "nginx",
		Tag:        "latest",
		Networks:   []*dockertest.Network{testNetwork(t)},
		Mounts: []string{
			fmt.Sprintf("%s:/etc/nginx/nginx.conf:ro", confp),
			fmt.Sprintf("%s:/etc/nginx/njs/sleep.js:ro", sp),
		},
	}
	e, err := pool.RunWithOptions(opt)
	if err != nil {
		t.Fatalf("Could not start resource: %s", err)
	}
	t.Cleanup(func() {
		if err := pool.Purge(e); err != nil {
			t.Fatalf("Could not purge resource: %s", err)
		}
	})

	var urlstr string
	if err := pool.Retry(func() error {
		urlstr = fmt.Sprintf("http://localhost:%s/", e.GetPort("80/tcp"))
		u, err := url.Parse(urlstr)
		if err != nil {
			return err
		}
		if _, err := http.Get(u.String()); err != nil {
			return err
		}
		return nil
	}); err != nil {
		t.Fatalf("Could not connect to NGINX server: %s", err)
	}
	return urlstr
}

func testNetwork(t testing.TB) *dockertest.Network {
	t.Helper()
	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("Could not connect to docker: %s", err)
	}
	ns, err := pool.NetworksByName(networkName)
	if err != nil {
		t.Fatalf("Could not connect to docker: %s", err)
	}
	switch len(ns) {
	case 0:
		n, err := pool.CreateNetwork(networkName)
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() {
			_ = pool.RemoveNetwork(n)
		})
		return n
	case 1:
		// deletion of network is left to the function that created it.
		return &ns[0]
	default:
		t.Fatalf("Could not connect to docker: %s", err)
	}
	return nil
}
