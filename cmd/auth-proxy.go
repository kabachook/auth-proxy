package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	flag.Parse()

	config, err := config.Load(configFile)
	if err != nil {
		log.Fatal(err)
	}

	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()

	sugar.Infof("Loaded config: %+v", config)

	proxy := proxy.New(*config, *logger)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr: config.Listen,
		Handler: proxy.Handler(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			sugar.Fatalw("Server listen error", err)
		}
	}()
	sugar.Infof("Server started at %s", config.Listen)

	<-done
	sugar.Info("Server stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		proxy.Shutdown()
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		sugar.Fatalf("Server Shutdown Failed:%+v", err)
	}
	sugar.Info("Server Exited Properly")
}
