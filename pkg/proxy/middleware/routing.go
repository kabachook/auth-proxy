package middleware

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kabachook/auth-proxy/pkg/config"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

// RoutingMiddleware provides routing features, i.e. routes requests to backends
func NewRoutingMiddleware(routes []config.Route, backends map[string]config.Backend, logger zap.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		l := logger.Named("routing")
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			host := strings.Split(r.Host, ":")

			// TODO: impove backend search, probably make Backend pointer in Route
			found := false
			for _, route := range routes {
				if route.Match.Host == "*" || route.Match.Host == host[0] {
					backend, ok := backends[route.Backend]
					if !ok {
						logger.Sugar().Errorf("ERROR: Can't find backend %s for host %s", route.Backend, route.Match.Host)
						w.WriteHeader(http.StatusBadRequest)
						return
					}
					r.URL.Scheme = backend.Scheme
					r.URL.Host = fmt.Sprintf("%s:%d", backend.Host, backend.Port)
					r.Host = r.URL.Host

					l.Sugar().Debugf("found %s -> %s\n", host[0], r.URL)
					found = true
				}
				if found {
					break
				}
			}

			if !found {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("Host not found"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
