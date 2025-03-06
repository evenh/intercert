package logging

import (
	"sync"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	once      sync.Once
	zapLogger *zap.Logger
)

func ConfigureLogging(isDebug bool) *zap.Logger {
	once.Do(func() {
		level := zap.NewAtomicLevelAt(zap.InfoLevel)
		if isDebug {
			level = zap.NewAtomicLevelAt(zap.DebugLevel)
		}

		// Custom encoder
		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.TimeKey = "timestamp"
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encoderConfig.MessageKey = "message"
		encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

		c := zap.Config{
			Level:         level,
			Encoding:      "json",
			EncoderConfig: encoderConfig,
			OutputPaths:   []string{"stdout"},
			DisableCaller: true,
		}

		l, err := c.Build()
		defer l.Sync()
		if err != nil {
			panic(errors.Wrap(err, "failed to build logger"))
		}

		zapLogger = l
	})

	return zapLogger
}
