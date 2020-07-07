package middleware

import (
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// LoggingMiddleware logs requests
func NewLoggingMiddleware(logger zap.Logger, identityField string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			lrw := &loggingResponseWriter{w, http.StatusOK}
			next.ServeHTTP(lrw, r)

			username := r.Context().Value(identityField)

			//var baseFields = []interface{}{"ip", r.RemoteAddr, "method", r.Method, "path", r.URL.EscapedPath(), "code", lrw.statusCode}
			var baseFields = []zap.Field{
				zap.String("ip", r.RemoteAddr),
				zap.String("method", r.Method),
				zap.String("path", r.URL.EscapedPath()),
				zap.Int("code", lrw.statusCode),
			}
			if username != nil {
				logger.Info("", append(baseFields, zap.String(identityField, username.(string)))...)
			} else {
				logger.Info("", baseFields...)
			}
		})
	}
}
