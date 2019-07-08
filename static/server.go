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
	s.monitoringRoutes()

	return s.router.ServeHTTP
}
