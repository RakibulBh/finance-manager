package middleware

import (
	"net/http"
	"time"

	"github.com/rakibulbh/ai-finance-manager/internal/logger"
	"go.uber.org/zap"
)

type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.wroteHeader {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.ResponseWriter.Write(b)
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := wrapResponseWriter(w)

		next.ServeHTTP(wrapped, r)

		duration := time.Since(start)

		// Log to logs/http.log via logger.HTTPLog
		// "HTTP Request" is the message, fields contain details
		fields := []zap.Field{
			zap.String("method", r.Method),
			zap.String("url", r.URL.String()),
			zap.String("remote_addr", r.RemoteAddr),
			zap.Int("status", wrapped.Status()),
			zap.Duration("duration", duration),
			zap.String("user_agent", r.UserAgent()),
		}

		if wrapped.Status() >= 500 {
			// Also log to server.log if it's a server error, for visibility
			logger.Error("HTTP Server Error",
				zap.String("method", r.Method),
				zap.String("url", r.URL.String()),
				zap.Int("status", wrapped.Status()),
			)
		}

		logger.InfoHTTP("HTTP Request", fields...)
	})
}
