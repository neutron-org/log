package logger

import (
	"fmt"
	"io/ioutil"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// loggerCfgPath is the path to a logger configuration file. Expected to be either json or yaml.
	loggerCfgPath = os.Getenv("LOGGER_CFG_PATH")
	// loggerLevel defines the logger logging level.
	loggerLevel = os.Getenv("LOGGER_LEVEL")
)

// NewForContext returns a new zap logger with an injected "context" field that is always
// presented in the logger messages. If a LOGGER_CFG_PATH env variable is defined, the logger
// config is build using the given file. Otherwise, a default simple config is used.
func NewForContext(context string) (*zap.Logger, error) {
	cfg, err := loadCfg()
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

// loadCfg loads logger configuration from the given config file or, if it's empty, returns the
// default config.
func loadCfg() (*zap.Config, error) {
	if loggerCfgPath == "" {
		return defaultLoggerCfg(), nil
	}
	cfg, err := cfgFromFile()
	if err != nil {
		return nil, fmt.Errorf("failed to build logger config from provided file %s: %w", loggerCfgPath, err)
	}
	return cfg, nil
}

// defaultLoggerCfg creates a default logger config.
func defaultLoggerCfg() *zap.Config {
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(parseLevel())
	cfg.EncoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder
	return &cfg
}

// cfgFromFile loads logger configuration from loggerCfgPath.
func cfgFromFile() (*zap.Config, error) {
	unmarshalFn, err := getCfgUnmarshaller(loggerCfgPath)
	if err != nil {
		return nil, err
	}

	cfgData, err := ioutil.ReadFile(loggerCfgPath)
	if err != nil {
		return nil, fmt.Errorf("read cfg file error: %w", err)
	}
	cfg, err := unmarshalFn(cfgData)
	if err != nil {
		return nil, fmt.Errorf("unmarshal error: %w", err)
	}
	return cfg, nil
}

// parseLevel parses logger level env into zapcore.Level. The env value is expected to contain one of
// the zapcore.Level const values. On an unexpected value, an Info level is returned.
func parseLevel() zapcore.Level {
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(loggerLevel)); err != nil {
		return zapcore.InfoLevel
	}
	return level
}
