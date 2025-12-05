package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	encodingJSON    = "json"
	encodingConsole = "console"

	//stdout = "stdout"
	stderr = "stderr"

	keyTime       = "time"
	keyLevel      = "level"
	keyName       = "name"
	keyCaller     = "caller"
	keyMsg        = "logMsg"
	keyStacktrace = "stacktrace"
)

type Logger struct {
	*zap.SugaredLogger
}

var l *Logger

func init() {
	l = New()
}

func Log() *zap.SugaredLogger {
	return l.SugaredLogger
}

// New initialize new logger instance
func New() *Logger {
	logger := &Logger{}
	zapLogger := logger.newZapLogger()
	logger.SugaredLogger = zapLogger.Sugar()
	return logger
}

func (logger *Logger) newZapLogger() *zap.Logger {
	zapConfig := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.DebugLevel),
		Development:       true,
		DisableCaller:     false,
		DisableStacktrace: false,
		Encoding:          encodingConsole,
		EncoderConfig:     logger.newZapEncoderConfig(),
		OutputPaths:       []string{stderr, "log.log"},
		ErrorOutputPaths:  []string{stderr, "log.log"},
		InitialFields:     nil,
	}

	zapOptions := []zap.Option{zap.AddCallerSkip(1)}
	zapLogger, _ := zapConfig.Build(zapOptions...)
	return zapLogger
}

func (logger *Logger) newZapEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		MessageKey:     keyMsg,
		LevelKey:       keyLevel,
		TimeKey:        keyTime,
		NameKey:        keyName,
		CallerKey:      keyCaller,
		StacktraceKey:  keyStacktrace,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseColorLevelEncoder,
		EncodeTime:     zapcore.RFC3339TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
}
