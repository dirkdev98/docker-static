package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/NYTimes/gziphandler"

	"github.com/dirkdev98/docker-static/static"
)

var (
	portPtr        = flag.Int("port", 3000, "Port to listen on")
	path           = flag.String("path", "/public", "Path on which files will be found")
	enableFallback = flag.Bool("fallback", true, "Automatically try to serve index.html if file is not found")
	basicAuth      = flag.String("auth", "", "Basic authorization in form of username:password")
	maxAge         = flag.Int("maxAge", 3600, "Cache-Control header value")
)

func main() {
	flag.Parse()

	port := fmt.Sprintf(":%d", *portPtr)

	opts := &static.Options{
		Path:      *path,
		Fallback:  *enableFallback,
		BasicAuth: *basicAuth,
		MaxAge:    *maxAge,
	}

	staticHandler := gziphandler.GzipHandler(http.HandlerFunc(static.ServerHandler(opts)))

	fmt.Printf("{\"level\": \"info\", \"timestamp\": \"%s\", \"type\": \"HTTP_STATIC_START\", \"message\": \"Listening at 0.0.0.0%v for static files\"}", time.Now().Format(time.RFC3339), port)
	fmt.Println()

	err := http.ListenAndServe(port, staticHandler)
	if err != nil {
		fmt.Printf("{\"level\": \"error\", \"timestamp\": \"%s\", \"type\": \"HTTP_STATIC_ERROR\", \"message\": \"%s\"}", time.Now().Format(time.RFC3339), err)
		fmt.Println()
	}
}
