package log

import (
	"sync/atomic"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var _globalLogger = defaultLogger()

func defaultLogger() atomic.Value {
	encCfg := zapcore.EncoderConfig{
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      "C",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	cfg := zap.Config{
		DisableCaller:    true,
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:      true,
		Encoding:         "console",
		EncoderConfig:    encCfg,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, _ := cfg.Build()

	var v atomic.Value
	v.Store(logger)

	return v
}

func L() *zap.Logger {
	return _globalLogger.Load().(*zap.Logger)
}

func S() *zap.SugaredLogger {
	return L().Sugar()
}

func LWithSvcName(name string) *zap.Logger {
	return L().With(zap.String("Service", name))
}

func AgentUID(uid uint64) zap.Field {
	return zap.Uint64("AgentUID", uid)
}
