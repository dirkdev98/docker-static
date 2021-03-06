package static

import (
	"net/http"
)

func (s *server) handleHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		errorLogger(w.Write([]byte("OK")))
	}
}
