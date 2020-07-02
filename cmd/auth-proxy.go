package main

import (
	"log"
	"net/http"

	"github.com/kabachook/auth-proxy/pkg/config"
	"github.com/kabachook/auth-proxy/pkg/proxy"
	flag "github.com/spf13/pflag"
)

var (
	configFile string
)

func init() {
	flag.StringVarP(&configFile, "config", "c", "config.yaml", "Config filename")
}

func main() {
	flag.Parse()
	log.Printf("config: %s\n", configFile)

	config, err := config.Load(configFile)
	if err != nil {
		log.Fatal(err)
	}

	proxy := proxy.NewProxy(*config)
	log.Printf("proxy: %v", proxy)
	log.Printf("Started at %s", config.Listen)
	log.Fatal(http.ListenAndServe(config.Listen, proxy))
}
