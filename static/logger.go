package static

import (
	"encoding/json"
	"net/http"
	"os"
	"time"
)

type logData struct {
	Host           string
	Remote         string
	Method         string
	Path           string
	Status         int
	ResponseLength int
	UserAgent      string
	Duration       float64
}

type logWriter struct {
	http.ResponseWriter
	status int
	length int
	start  time.Time
}

type logWriterFunc func(w *logWriter, r *http.Request)

func (w *logWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *logWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = 200
	}
	n, err := w.ResponseWriter.Write(b)
	w.length += n
	return n, err
}

func (s *server) logMiddleware(h logWriterFunc) http.HandlerFunc {
	out := json.NewEncoder(os.Stdout)

	return func(w http.ResponseWriter, r *http.Request) {
		l := &logWriter{
			ResponseWriter: w,
			start:          time.Now(),
		}

		defer func() {
			err := out.Encode(logData{
				Host:           r.Host,
				Remote:         r.RemoteAddr,
				Method:         r.Method,
				Path:           r.URL.Path,
				Status:         l.status,
				ResponseLength: l.length,
				UserAgent:      r.UserAgent(),
				Duration:       time.Now().Sub(l.start).Seconds(),
			})

			if err != nil {
				errorLogger(0, err)
			}
		}()

		h(l, r)
	}
}
