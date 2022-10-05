package logger

import (
	"encoding/json"
	"fmt"
	"strings"

	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

// getCfgUnmarshaller returns a cfgUnmarshaller based on the provided cfg file extension.
func getCfgUnmarshaller(cfgFile string) (cfgUnmarshaller, error) {
	lowerCfgFile := strings.ToLower(cfgFile)
	switch {
	case strings.HasSuffix(lowerCfgFile, "json"):
		return jsonCfgUnmarshaller, nil
	case strings.HasSuffix(lowerCfgFile, "yaml"), strings.HasSuffix(lowerCfgFile, "yml"):
		return yamlCfgUnmarshaller, nil
	default:
		return nil, fmt.Errorf("unexpected format of cfg file")
	}
}

// cfgUnmarshaller is a type defining an algorithm of unmarshalling a logger cfg file into a
// *zap.Config instance.
type cfgUnmarshaller func(data []byte) (*zap.Config, error)

// jsonCfgUnmarshaller is a cfgUnmarshaller for json files.
var jsonCfgUnmarshaller = func(data []byte) (*zap.Config, error) {
	var cfg zap.Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// yamlCfgUnmarshaller is a cfgUnmarshaller for yaml files.
var yamlCfgUnmarshaller = func(data []byte) (*zap.Config, error) {
	var cfg zap.Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
