package middleware

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/kabachook/auth-proxy/pkg/config"
	"go.uber.org/zap"
	"net/http"
)

// AuthnMiddleware gets user from context and sets identity field for next middlewares
func NewAuthnMiddleware(authnConfig config.AuthnConfig, proxyConfig config.Proxy, identityField string, logger zap.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		l := logger.Named("authn")
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := r.Context().Value("user")
			l.Sugar().Debugf("got token: %+v", user.(*jwt.Token))

			username, ok := user.(*jwt.Token).Claims.(jwt.MapClaims)[authnConfig.JWT.Field]
			if !ok {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Error getting username"))
				return
			}

			r.Header.Add(proxyConfig.Header, fmt.Sprint(username))
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), identityField, username)))
		})
	}
}
