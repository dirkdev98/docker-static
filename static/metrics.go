package static

import (
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"strconv"
	"time"
)

const (
	prometheusPrefix = "docker_static_"
)

var (
	authCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: prometheusPrefix + "auth_count",
		Help: "The number of authorizations by success and failure",
	}, []string{"type"})

	requestCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: prometheusPrefix + "request_count",
		Help: "The number of requests per request path, request method and response status",
	}, []string{"path", "method", "status"})

	responseTime = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    prometheusPrefix + "response_time",
		Help:    "The history of response time by request path",
		Buckets: []float64{0.005, 0.01, 0.05, 0.25, 0.5, 1, 2.5, 5},
	}, []string{"path"})
)

func InitMetrics() {
	prometheus.MustRegister(authCount, requestCount, responseTime)
}

func (s *server) wrapWithMetrics(h logWriterFunc) logWriterFunc {
	return func(w *logWriter, r *http.Request) {
		path := r.URL.Path
		method := r.Method

		defer func() {
			requestCount.WithLabelValues(path, method, strconv.FormatInt(int64(w.status), 10)).Inc()
			responseTime.WithLabelValues(path).Observe(time.Now().Sub(w.start).Seconds())
		}()

		h(w, r)
	}
}
