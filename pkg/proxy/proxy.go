package proxy

import (
	"github.com/kabachook/auth-proxy/pkg/config"
)

// Proxy : auth-proxy struct
type Proxy struct {
	cfg config.Config
}

// NewProxy : creates new proxy
func NewProxy(cfg config.Config) *Proxy {
	return &Proxy{
		cfg: cfg,
	}
}
