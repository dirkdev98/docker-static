package main

import (
	"github.com/dirkdev98/docker-static/static"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCacheControl(t *testing.T) {
	opts := &static.Options{
		Path:         ".",
		FallbackPath: "index.html",
		BasicAuth:    "",
		MaxAge:       125,
	}

	server := httptest.NewServer(static.ServerHandler(opts))

	t.Run("Use provided cache control", func(t *testing.T) {
		res, err := http.Get(server.URL)
		if err != nil {
			log.Fatalln(err)
		}

		if res.StatusCode != 200 {
			t.Fail()
		}
		if res.Header.Get("Cache-Control") != `max-age=125` {
			t.Fail()
		}
	})

	server.Close()
}

func TestStaticAuth(t *testing.T) {
	opts := &static.Options{
		Path:         ".",
		FallbackPath: "index.html",
		BasicAuth:    "user:pass",
		MaxAge:       10,
	}

	server := httptest.NewServer(static.ServerHandler(opts))

	t.Run("Fail without authorization", func(t *testing.T) {

		res, err := http.Get(server.URL)
		if err != nil {
			log.Fatalln(err)
		}

		if res.StatusCode != 401 {
			t.Fail()
		}
		if res.Header.Get("WWW-Authenticate") != `Basic realm="Restricted"` {
			t.Fail()
		}
	})

	t.Run("Success with authorization", func(t *testing.T) {

		client := &http.Client{}
		req, err := http.NewRequest("GET", server.URL, nil)
		if err != nil {
			log.Fatalln(err)
		}

		req.SetBasicAuth("user", "pass")
		res, err := client.Do(req)

		if err != nil {
			log.Fatalln(err)
		}

		if res.StatusCode != 200 {
			t.Fail()
		}
	})

	t.Run("Error with false authorization", func(t *testing.T) {
		client := &http.Client{}
		req, err := http.NewRequest("GET", server.URL, nil)
		if err != nil {
			log.Fatalln(err)
		}

		req.SetBasicAuth("user", "pas")
		res, err := client.Do(req)

		if err != nil {
			log.Fatalln(err)
		}

		if res.StatusCode != 401 {
			t.Fail()
		}
	})

	server.Close()

	opts.BasicAuth = ""
	server = httptest.NewServer(static.ServerHandler(opts))

	t.Run("Success without authorization", func(t *testing.T) {
		res, err := http.Get(server.URL)
		if err != nil {
			log.Fatalln(err)
		}

		if res.StatusCode != 200 {
			t.Fail()
		}
	})

	server.Close()
}

func TestPathAndFallback(t *testing.T) {
	opts := &static.Options{
		Path:         ".",
		FallbackPath: "",
		BasicAuth:    "",
		MaxAge:       10,
	}

	server := httptest.NewServer(static.ServerHandler(opts))

	t.Run("GET /", func(t *testing.T) {
		res, err := http.Get(server.URL + "/")
		if err != nil {
			log.Fatalln(err)
		}

		if res.StatusCode != 200 {
			t.Fail()
		}
	})

	t.Run("GET /index.html", func(t *testing.T) {
		res, err := http.Get(server.URL + "/index.html")
		if err != nil {
			log.Fatalln(err)
		}

		if res.StatusCode != 200 {
			t.Fail()
		}
	})

	t.Run("GET /file.html", func(t *testing.T) {
		res, err := http.Get(server.URL + "/file.html")
		if err != nil {
			log.Fatalln(err)
		}

		if res.StatusCode != 404 {
			t.Fail()
		}
	})

	server.Close()
	opts.FallbackPath = "index.html"

	server = httptest.NewServer(static.ServerHandler(opts))

	t.Run("GET / - with fallback", func(t *testing.T) {
		res, err := http.Get(server.URL + "/")
		if err != nil {
			log.Fatalln(err)
		}

		if res.StatusCode != 200 {
			t.Fail()
		}
	})

	t.Run("GET /index.html - with fallback", func(t *testing.T) {
		res, err := http.Get(server.URL + "/index.html")
		if err != nil {
			log.Fatalln(err)
		}

		if res.StatusCode != 200 {
			t.Fail()
		}
	})

	t.Run("GET /file.html - with fallback", func(t *testing.T) {
		res, err := http.Get(server.URL + "/file.html")
		if err != nil {
			log.Fatalln(err)
		}

		if res.StatusCode != 200 {
			t.Fail()
		}
	})

	server.Close()
}

func TestMonitoring(t *testing.T) {
	opts := &static.Options{
		Path:         ".",
		FallbackPath: "",
		BasicAuth:    "",
		MaxAge:       10,
	}

	server := httptest.NewServer(static.MonitoringHandler(opts))

	t.Run("GET /health always ok", func(t *testing.T) {
		res, err := http.Get(server.URL + "/health")
		if err != nil {
			log.Fatalln(err)
		}

		if res.StatusCode != 200 {
			t.Fail()
		}
	})

	t.Run("GET /metrics", func(t *testing.T) {
		res, err := http.Get(server.URL + "/metrics")
		if err != nil {
			log.Fatalln(err)
		}

		if res.StatusCode != 200 {
			t.Fail()
		}
	})

	server.Close()

	t.Run("Init metrics does not panic", func(t *testing.T) {
		static.InitMetrics()
	})
}