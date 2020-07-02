package proxy

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"

	"github.com/kabachook/auth-proxy/pkg/config"
)

// Proxy : auth-proxy struct
type Proxy struct {
	cfg      config.Config
	backends map[string]config.Backend
	handler  http.Handler
}

func base(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user")

		username, ok := user.(*jwt.Token).Claims.(jwt.MapClaims)["username"] // TODO: find out how to unhardcode it
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error getting username"))
			return
		}

		r.Header.Add("X-Username", fmt.Sprint(username))
		next.ServeHTTP(w, r)
	})
}

// New : creates new proxy
func New(cfg config.Config) *Proxy {
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.Authn.JWT.Secret), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})

	reverseProxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
		},
	}
	handler := base(reverseProxy)
	handler = handlers.CombinedLoggingHandler(os.Stdout, handler)
	handler = jwtMiddleware.Handler(handler)

	return &Proxy{
		cfg:      cfg,
		backends: config.BackendsToMap(cfg.Backends),
		handler:  handler,
	}
}

// TODO: move to middleware
func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	host := strings.Split(r.Host, ":")

	// TODO: impove backend search, probably make Backend pointer in Route
	found := false
	for _, route := range p.cfg.Routes {
		if route.Match.Host == "*" || route.Match.Host == host[0] {
			backend, ok := p.backends[route.Backend]
			if !ok {
				log.Printf("ERROR: Can't find backend %s for host %s", route.Backend, route.Match.Host)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			r.URL.Scheme = backend.Scheme
			r.URL.Host = fmt.Sprintf("%s:%d", backend.Host, backend.Port)
			r.Host = r.URL.Host // Make ReverseProxy use req.URL

			log.Printf("Found %s -> %s\n", host[0], r.URL)
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

	p.handler.ServeHTTP(w, r)
}
