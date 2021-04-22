package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func BuildLoggerOrDie(debug bool) *zap.Logger {
	logger, err := buildLogger(debug)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not build logger, exiting: %v\n", err)
		os.Exit(1)
	}
	zap.ReplaceGlobals(logger)
	return logger
}

func BuildLogger(debug bool) (*zap.Logger, error) {
	logger, err := buildLogger(debug)
	zap.ReplaceGlobals(logger)
	return logger, err
}

func buildLogger(debug bool) (*zap.Logger, error) {
	var err error
	var logger *zap.Logger
	zapLevel := zapcore.InfoLevel
	if debug {
		zapLevel = zapcore.DebugLevel
	}
	logger, err = zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(zapLevel),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:   "message",
			LevelKey:     "level",
			EncodeLevel:  zapcore.LowercaseLevelEncoder,
			TimeKey:      "time",
			EncodeTime:   zapcore.EpochTimeEncoder,
			CallerKey:    "caller",
			EncodeCaller: zapcore.FullCallerEncoder,
		},
	}.Build()
	return logger, err
}
