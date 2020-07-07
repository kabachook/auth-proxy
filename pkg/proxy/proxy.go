package proxy

import (
	"context"
	"fmt"
	"github.com/kabachook/auth-proxy/pkg/acl"
	"github.com/kabachook/auth-proxy/pkg/authz"
	"net/http"
	"net/http/httputil"
	"strings"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/kabachook/auth-proxy/pkg/config"
)

var (
	JWTIdentityField = "username"
)

// Proxy : auth-proxy struct
type Proxy struct {
	cfg      config.Config
	backends map[string]config.Backend
	handler  http.Handler
	logger   zap.Logger
	authz    authz.Authz
}

// loggingMiddleware logs requests, depending on a position in middleware chain
func loggingMiddleware(cfg config.AuthnConfig, logger zap.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			username := r.Context().Value(JWTIdentityField).(string)

			logger.Sugar().Infow("request", "host", r.Host, "url", r.URL.EscapedPath(), cfg.JWT.Field, username)
			next.ServeHTTP(w, r)
		})
	}
}

// authnMiddleware gets user from context and sets identity field for following middlewares
func authnMiddleware(authnConfig config.AuthnConfig, proxyConfig config.Proxy, logger zap.Logger) mux.MiddlewareFunc {
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
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), JWTIdentityField, username)))
		})
	}
}

// authzMiddleware authenticates request based on identity field from context using authz
func authzMiddleware(authz authz.Authz, logger zap.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		l := logger.Named("authz")
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			username := r.Context().Value(JWTIdentityField).(string)
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

// aclMiddleware authenticates request based on identity field and path access control list
func aclMiddleware(acl acl.ACL, logger zap.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		l := logger.Named("acl")
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			username := r.Context().Value(JWTIdentityField).(string)
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

// routingMiddleware provides routing features, i.e. routes requests to backends
func routingMiddleware(routes []config.Route, backends map[string]config.Backend, logger zap.Logger) mux.MiddlewareFunc {
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

					l.Sugar().Infof("found %s -> %s\n", host[0], r.URL)
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

// New : creates new proxy
func New(cfg config.Config, logger zap.Logger) *Proxy {
	router := mux.NewRouter()
	reverseProxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
		},
	}
	router.PathPrefix("/").Handler(reverseProxy)

	ldapAuthz, err := authz.NewLDAPAuthz(cfg.Authz.LDAP)
	if err != nil {
		logger.Sugar().Fatalw("Error creating LDAPAuthz", err.Error())
	}
	err = ldapAuthz.Open()
	if err != nil {
		logger.Sugar().Fatalw("Can't open LDAP connection", err.Error())
	}

	middlewares := []mux.MiddlewareFunc{
		jwtmiddleware.New(jwtmiddleware.Options{
			ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
				return []byte(cfg.Authn.JWT.Secret), nil
			},
			SigningMethod: jwt.SigningMethodHS256,
		}).Handler,
		authnMiddleware(cfg.Authn, cfg.Proxy, logger),
		authzMiddleware(ldapAuthz, logger),
		aclMiddleware(acl.NewSimpleACL(cfg.Authz.ACL), logger),
		routingMiddleware(cfg.Routes, config.BackendsToMap(cfg.Backends), logger),
		loggingMiddleware(cfg.Authn, logger),
	}

	router.Use(middlewares...)

	return &Proxy{
		cfg:      cfg,
		backends: config.BackendsToMap(cfg.Backends),
		handler:  router,
		logger:   logger,
		authz:    ldapAuthz,
	}
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			p.logger.Sugar().Error(err)
		}
	}()

	p.handler.ServeHTTP(w, r)
}

func (p *Proxy) Handler() http.Handler {
	return p.handler
}

func (p *Proxy) Shutdown() {
	p.logger.Sugar().Info("Closing connections")
	p.authz.Close()
}
