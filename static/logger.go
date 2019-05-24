package static

import (
	"encoding/json"
	"net/http"
	"os"
	"time"
)

const (
  level = "info"
  logType = "HTTP_STATIC_LOG"
)

type logData struct {
    Level          string  `json:"level"`
    Timestamp      string  `json:"timestamp"`
    Type           string  `json:"type"`
	Host           string  `json:"host"`
	Remote         string  `json:"remote"`
	Method         string  `json:"method"`
	Path           string  `json:"path"`
	Status         int     `json:"status"`
	ResponseLength int     `json:"responseLength"`
	UserAgent      string  `json:"userAgent"`
	Duration       float64 `json:"duration"`
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
			    Level:          level,
			    Type:           logType,
			    Timestamp:      time.Now().Format(time.RFC3339),
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
