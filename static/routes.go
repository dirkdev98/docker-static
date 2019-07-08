package static

func (s *server) routes() {
	s.router.HandleFunc("/", s.logMiddleware(s.authMiddleware(s.handleStatic())))
}

func (s *server) monitoringRoutes() {
	s.router.HandleFunc("/health", s.handleHealth())
}
