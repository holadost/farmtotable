package util

import (
	"encoding/json"
	"go.uber.org/zap"
	"os"
	"strings"
)

func NewJSONLogger() *zap.Logger {
	// TODO: Get the correct log level here.
	rawJSON := []byte(`{
	  "level": "info",
	  "encoding": "json",
	  "outputPaths": ["stdout", "/tmp/out.log"],
	  "errorOutputPaths": ["stderr", "/tmp/err.log"],
	  "encoderConfig": {
	    "messageKey": "message",
	    "levelKey": "level",
	    "levelEncoder": "lowercase"
	  }
	}`)

	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	return logger
}

func GetLogLevel() string {
	value, exists := os.LookupEnv("FTT_GLOBAL_LOG_LEVEL")
	if !exists {
		return "info"
	}
	return strings.ToLower(value)
}
