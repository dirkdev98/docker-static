package static

import (
	"encoding/base64"
	"log"
	"net/http"
	"strings"
)

func parseAuth(authString string) (string, string) {
	identity := strings.Split(authString, ":")
	if len(identity) != 2 {
		log.Fatalln("Provided basic auth flag but invalid value. See -help")
	}

	return identity[0], identity[1]
}

func (s *server) authMiddleware(next logWriterFunc) logWriterFunc {
	// No setup and early return if no auth is needed.
	if len(s.staticOptions.BasicAuth) == 0 {
		return next
	}

	username, password := parseAuth(s.staticOptions.BasicAuth)

	return func(w *logWriter, r *http.Request) {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)

		auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)

		if auth[0] != "Basic" {
			http.Error(w, "authorization failed", http.StatusUnauthorized)
			return
		}

		payload, err := base64.StdEncoding.DecodeString(auth[1])

		if err != nil {
			log.Println(err)
			http.Error(w, "authorization failed", http.StatusUnauthorized)
			return
		}

		pair := strings.SplitN(string(payload), ":", 2)
		if strings.Compare(pair[0], username) != 0 || strings.Compare(pair[1], password) != 0 {
			http.Error(w, "authorization failed", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
