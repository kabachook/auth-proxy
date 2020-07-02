package proxy

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"

	"github.com/gorilla/handlers"

	"github.com/kabachook/auth-proxy/pkg/config"
)

// Proxy : auth-proxy struct
type Proxy struct {
	cfg      config.Config
	backends map[string]config.Backend
	handler  http.Handler
}

// NewProxy : creates new proxy
func NewProxy(cfg config.Config) *Proxy {
	reverseProxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
		},
	}

	handler := handlers.CombinedLoggingHandler(os.Stdout, reverseProxy)

	return &Proxy{
		cfg:      cfg,
		backends: config.BackendsToMap(cfg.Backends),
		handler:  handler,
	}
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	host := strings.Split(req.Host, ":")

	// TODO: impove backend search, probably make Backend pointer in Route
	found := false
	for _, route := range p.cfg.Routes {
		if route.Match.Host == "*" || route.Match.Host == host[0] {
			backend, ok := p.backends[route.Backend]
			if !ok {
				log.Printf("ERROR: Can't find backend %s for host %s", route.Backend, route.Match.Host)
				w.WriteHeader(http.StatusNotFound)
				return
			}
			req.URL.Scheme = backend.Scheme
			req.URL.Host = fmt.Sprintf("%s:%d", backend.Host, backend.Port)
			req.Host = req.URL.Host // Make ReverseProxy use req.URL

			log.Printf("Found %s -> %s\n", host[0], req.URL)
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

	p.handler.ServeHTTP(w, req)
}
