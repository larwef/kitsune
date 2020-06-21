package main

import (
	"fmt"
	"github.com/larwef/kitsune/repository/memory"
	"github.com/larwef/kitsune/server"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"net/http"
	"strings"
	"time"
)

// Version injected at compile time.
var Version = "No version provided"

func main() {
	setupLogger()
	defer zap.L().Sync()
	setupConfig()

	zap.S().Infow("Staring application",
		"version", Version,
	)

	ks := server.NewServer(memory.NewRepository())

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", viper.GetInt("port")),
		Handler:      ks.GetRouter(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	zap.S().Infof("Starting server on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		zap.S().Errorw("Listen and serve returned an error",
			"error", err,
		)
	}

	zap.S().Info("Application exited")
}

func setupConfig() {
	viper.SetEnvPrefix("KITSUNE")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	viper.SetDefault("PORT", 8080)
}

func setupLogger() {
	zapConfig := zap.NewProductionEncoderConfig()
	zapConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logConfig := zap.NewProductionConfig()
	logConfig.EncoderConfig = zapConfig

	logger, err := logConfig.Build()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}

	zap.ReplaceGlobals(logger)
}
