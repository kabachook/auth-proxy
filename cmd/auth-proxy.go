package main

import (
	"log"
	"net/http"

	"github.com/kabachook/auth-proxy/pkg/config"
	"github.com/kabachook/auth-proxy/pkg/proxy"
	flag "github.com/spf13/pflag"
	"go.uber.org/zap"
)

var (
	configFile string
)

func init() {
	flag.StringVarP(&configFile, "config", "c", "config.yaml", "Config filename")
}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()

	flag.Parse()

	config, err := config.Load(configFile)
	if err != nil {
		log.Fatal(err)
	}
	sugar.Infof("Loaded config: %+v", config)

	proxy := proxy.New(*config, *logger)
	sugar.Infof("Started at %s", config.Listen)
	sugar.Fatal(http.ListenAndServe(config.Listen, proxy))
}
