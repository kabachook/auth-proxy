package proxy

import (
	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/kabachook/auth-proxy/pkg/acl"
	"github.com/kabachook/auth-proxy/pkg/authz"
	"github.com/kabachook/auth-proxy/pkg/proxy/middleware"
	"go.uber.org/zap"
	"net/http"
	"net/http/httputil"

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
		middleware.NewLoggingMiddleware(logger, JWTIdentityField),
		jwtmiddleware.New(jwtmiddleware.Options{
			ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
				return []byte(cfg.Authn.JWT.Secret), nil
			},
			SigningMethod: jwt.SigningMethodHS256,
		}).Handler,
		middleware.NewAuthnMiddleware(cfg.Authn, cfg.Proxy, JWTIdentityField, logger),
		middleware.NewAuthzMiddleware(ldapAuthz, JWTIdentityField, logger),
		middleware.NewACLMiddleware(acl.NewSimpleACL(cfg.Authz.ACL), JWTIdentityField, logger),
		middleware.NewRoutingMiddleware(cfg.Routes, config.BackendsToMap(cfg.Backends), logger),
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
