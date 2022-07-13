package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// statusWriter is custom http.ResponseWriter that captures status and size of response
type statusWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = 200
	}
	n, err := w.ResponseWriter.Write(b)
	w.size += n
	return n, err
}

func Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := statusWriter{ResponseWriter: w}

		next(&sw, r)

		fields := []zapcore.Field{
			zap.Int("status", sw.status),
			zap.Int("response_size", sw.size),
			zap.String("latency", time.Since(start).String()),
			zap.String("method", r.Method),
			zap.String("uri", r.RequestURI),
			zap.String("host", r.Host),
			zap.String("remote_ip", r.RemoteAddr),
		}

		n := sw.status
		switch {
		case n >= 500:
			zap.L().Error("Server error", fields...)
		case n >= 400:
			zap.L().Warn("Client error", fields...)
		case n >= 300:
			zap.L().Info("Redirection", fields...)
		default:
			zap.L().Info("Success", fields...)
		}
	}
}
