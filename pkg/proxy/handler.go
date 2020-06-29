/*
This file implements proxy handler, which routes the request
to specified backend based on Host header
*/
package proxy

import (
	"net/http"

	"github.com/kabachook/auth-proxy/pkg/config"
)

type ProxyHandler struct {
	baseHandler http.Handler
	backends    map[string]config.Backend
}

func NewProxyHandler(handler http.Handler, backends map[string]config.Backend) *ProxyHandler {
	return &ProxyHandler{
		baseHandler: handler,
		backends:    backends,
	}
}

func (p *ProxyHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// TODO: implement Host switching logic
	p.baseHandler.ServeHTTP(w, req)
}
