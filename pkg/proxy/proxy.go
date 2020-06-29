package proxy

import (
	"net/http"
	"net/http/httputil"

	"github.com/kabachook/auth-proxy/pkg/config"
)

// Proxy : auth-proxy struct
type Proxy struct {
	cfg     config.Config
	handler http.Handler
}

// NewProxy : creates new proxy
func NewProxy(cfg config.Config) *Proxy {
	reverseProxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {},
	}

	handler := NewProxyHandler(reverseProxy, config.BackendsToMap(cfg.Backends))

	return &Proxy{
		cfg:     cfg,
		handler: handler,
	}
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	p.handler.ServeHTTP(w, req)
}
