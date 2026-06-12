package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Log is the global logger instance.
var Log *zap.Logger

// Init initializes the Uber Zap logger.
func Init() {
	config := zap.NewDevelopmentConfig()
	// Enable colorized level encoder for beautiful dev console log output
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	var err error
	Log, err = config.Build()
	if err != nil {
		panic("Failed to initialize Zap logger: " + err.Error())
	}
}
