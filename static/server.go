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
	s.monitoringRoutes()
	s.routes()

	return s.router.ServeHTTP
}
