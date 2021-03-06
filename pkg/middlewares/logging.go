package middlewares

import (
	"net/http"
	"time"

	"github.com/spy16/droplets/pkg/logger"
)

// WithRequestLogging adds logging to the given handler. Every request handled by
// 'next' will be logged with request information such as path, method, latency,
// client-ip, response status code etc. Logging will be done at info level only.
// Also, injects a logger into the ResponseWriter which can be later used by the
// handlers to perform additional logging.
func WithRequestLogging(logger logger.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
		wrappedWr := &wrappedWriter{
			ResponseWriter: wr,
			Logger:         logger,
			wroteStatus:    http.StatusOK,
		}

		start := time.Now()
		next.ServeHTTP(wrappedWr, req)
		duration := time.Now().Sub(start)

		info := map[string]interface{}{
			"latency": duration,
			"status":  wrappedWr.wroteStatus,
		}

		logger.
			WithFields(requestInfo(req)).
			WithFields(info).
			Infof("request completed with code %d", wrappedWr.wroteStatus)
	})
}
