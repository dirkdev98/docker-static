package main

import (
	"flag"
	"fmt"
	"github.com/dirkdev98/docker-static/static"
	"log"
	"net/http"
)

var (
	portPtr           = flag.Int("port", 3000, "Port to listen on")
	enableMonitoring  = flag.Bool("monitoring", true, "Enable or disable monitoring")
	monitoringPortPtr = flag.Int("monitoringPort", 9600, "Monitoring port to listen on")
	path              = flag.String("path", "/public", "Path on which files will be found")
	fallbackPath      = flag.String("fallback", "", "Default file that will be served")
	basicAuth         = flag.String("auth", "", "Basic authorization in form of username:password")
	maxAge            = flag.Int("maxAge", 3600, "Cache-Control header value")
)

func main() {
	flag.Parse()

	port := fmt.Sprintf(":%d", *portPtr)
	monitoringPort := fmt.Sprintf(":%d", *monitoringPortPtr)

	static.InitMetrics()

	opts := &static.Options{
		Path:         *path,
		FallbackPath: *fallbackPath,
		BasicAuth:    *basicAuth,
		MaxAge:       *maxAge,
	}

	if *enableMonitoring {
		monitoringHandler := static.MonitoringHandler(opts)

		log.Printf("Listening at 0.0.0.0%v for monitoring", monitoringPort)
		if err := http.ListenAndServe(port, monitoringHandler); err != nil {
			log.Fatalln(err)
		}
	}

	staticHandler := static.ServerHandler(opts)

	log.Printf("Listening at 0.0.0.0%v for static files", port)
	if err := http.ListenAndServe(port, staticHandler); err != nil {
		log.Fatalln(err)
	}

}
