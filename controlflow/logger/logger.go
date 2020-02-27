package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

// NewLogger returns a new zap logger
func NewLogger(level string) *zap.Logger {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder // "2006-01-02T15:04:05.000Z0700"
	atom := zap.NewAtomicLevel()
	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		atom,
	))
	defer logger.Sync()
	var l zapcore.Level
	switch level {
	case "debug":
		l = zap.DebugLevel
	case "info":
		l = zap.InfoLevel
	case "warn":
		l = zap.WarnLevel
	case "error":
		l = zap.ErrorLevel
	case "fatal":
		l = zap.FatalLevel
	case "panic":
		l = zap.PanicLevel
	default:
		l = zap.DebugLevel
		logger.Warn("Log level '" + level + "' not recognized. Default set to Debug")
	}
	atom.SetLevel(l)
	return logger
}

// NewLoggerWithName ...
func NewLoggerWithName(name, level string, options ...zap.Option) (*zap.Logger, error) {
	atom := zap.NewAtomicLevel()
	cfg := zap.Config{
		Level:             atom,
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          "json",
		EncoderConfig:     zap.NewProductionEncoderConfig(),
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stderr"},
		InitialFields:     nil,
	}
	cfg.EncoderConfig.MessageKey = "message"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // "2006-01-02T15:04:05.000Z0700"
	logger, err := cfg.Build(options...)
	// err when build encoder and openSink (output)
	if err != nil {
		return nil, fmt.Errorf("build zap logger failed: %+v", err)
	}
	// name it
	logger = logger.Named(name)
	defer logger.Sync()
	var l zapcore.Level
	switch level {
	case "debug":
		l = zap.DebugLevel
	case "info":
		l = zap.InfoLevel
	case "warn":
		l = zap.WarnLevel
	case "error":
		l = zap.ErrorLevel
	case "fatal":
		l = zap.FatalLevel
	case "panic":
		l = zap.PanicLevel
	default:
		return nil, fmt.Errorf("log level '%s' not recognized, only support 'debug/info/warn/error/fatal/panic'", level)
	}
	atom.SetLevel(l)
	return logger, nil
}
