package middleware

import (
	"github.com/gorilla/mux"
	"github.com/kabachook/auth-proxy/pkg/authz"
	"go.uber.org/zap"
	"net/http"
)

// AuthzMiddleware authenticates request based on identity field from context using authz
func NewAuthzMiddleware(authz authz.Authz, identityField string, logger zap.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		l := logger.Named("authz")
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			username := r.Context().Value(identityField).(string)
			ok, err := authz.Authorize(username)

			l.Sugar().Debugf("%s: %t", username, ok)

			if err != nil {
				logger.Sugar().Errorw(err.Error())
				w.WriteHeader(http.StatusBadGateway)
				return
			}
			if !ok {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
