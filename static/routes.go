package static

import "github.com/prometheus/client_golang/prometheus/promhttp"

func (s *server) routes() {
	s.router.HandleFunc("/", s.logMiddleware(s.wrapWithMetrics(s.authMiddleware(s.handleStatic()))))
}

func (s *server) monitoringRoutes() {
	s.router.HandleFunc("/metrics", promhttp.Handler().ServeHTTP)
	s.router.HandleFunc("/health", s.handleHealth())
}
