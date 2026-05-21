package metrics

import (
	"net/http"
	"time"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func (m *Metrics) MiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rec := &statusRecorder{
			ResponseWriter: w,
			status:         200,
		}

		next.ServeHTTP(rec, r)

		m.RequestTotal.WithLabelValues(
			r.Method,
			r.URL.Path,
			http.StatusText(rec.status),
		).Inc()

		m.RequestDuration.WithLabelValues(
			r.Method,
			r.URL.Path,
		).Observe(time.Since(start).Seconds())
	})
}
