package static

import "net/http"

type server struct {
	router        *http.ServeMux
	staticOptions *Options
}

func ServerHandler(opts *Options) http.HandlerFunc {
	s := server{
		router:        http.NewServeMux(),
		staticOptions: opts,
	}
	s.routes()

	return s.router.ServeHTTP
}

func MonitoringHandler(opts *Options) http.HandlerFunc {
	s := server{
		router:        http.NewServeMux(),
		staticOptions: opts,
	}
	s.monitoringRoutes()

	return s.router.ServeHTTP
}
