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

	config, err := config.Load(configFile)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Loaded config: %+v\n", config)

	proxy := proxy.New(*config)
	log.Printf("Started at %s", config.Listen)
	log.Fatal(http.ListenAndServe(config.Listen, proxy))
}
