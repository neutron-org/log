package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logLevel = os.Getenv("LOG_LEVEL")

// NewForContext returns a new zap logger with an injected "context" field that is always
// presented in the logger messages.
func NewForContext(context string) (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder
	cfg.Level = zap.NewAtomicLevelAt(parseLevel())
	l, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	l = l.With(zap.String("context", context))
	return l, nil
}

// parseLevel parses env LOG_LEVEL into zapcore.Level. The env value is expected to contain one of
// the zapcore.Level const values. On an unexpected value, an Info level is returned.
func parseLevel() zapcore.Level {
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(logLevel)); err != nil {
		return zapcore.InfoLevel
	}
	return level
}
