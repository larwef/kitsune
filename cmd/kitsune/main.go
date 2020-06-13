package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
)

var Version = "No version provided"

func main() {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	logConfig := zap.NewProductionConfig()
	logConfig.EncoderConfig = config

	logger, err := logConfig.Build()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()

	zap.ReplaceGlobals(logger)

	zap.S().Infow("Staring application",
		"version", Version,
	)

	// Run application logic

	zap.S().Infow("Application exited")
}
