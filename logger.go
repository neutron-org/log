package logger

import (
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// loggerPrefix is the prefix for env variables that configure the logger.
const loggerPrefix = "LOGGER"

// NewForContext returns a new zap logger with an injected "context" field that is always
// presented in the logger messages. If a LOGGER_CFG_PATH env variable is defined, the logger
// config is build using the given file. Otherwise, a default simple config is used.
func NewForContext(context string) (*zap.Logger, error) {
	cfg, err := loadConfig()
	if err != nil {
		return nil, err
	}
	l, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	l = l.With(zap.String("context", context))
	return l, nil
}

// loadConfig initializes a default production logger config with parameters overwritten by
// corresponding env vars if set.
func loadConfig() (*zap.Config, error) {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder
	err := envconfig.Process(loggerPrefix, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
