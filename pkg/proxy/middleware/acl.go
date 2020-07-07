package middleware

import (
	"github.com/gorilla/mux"
	"github.com/kabachook/auth-proxy/pkg/acl"
	"go.uber.org/zap"
	"net/http"
)

// ACLMiddleware authenticates request based on identity field and path access control list
func NewACLMiddleware(acl acl.ACL, identityField string, logger zap.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		l := logger.Named("acl")
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			username := r.Context().Value(identityField).(string)
			pass := acl.Check(username, r.URL.Path)
			l.Sugar().Debugf("ACL: %t ->  %s %s", pass, username, r.URL.Path)

			if !pass {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
