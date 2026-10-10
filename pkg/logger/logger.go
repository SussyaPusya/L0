package logger

import (
	"go.uber.org/zap"
)

type Logger interface {
	Info(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Debug(msg string, fields ...zap.Field)
	Sync() error
}

type LoggerZap struct {
	zap *zap.Logger
}

func Newlogger() (*LoggerZap, error) {
	z, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	return &LoggerZap{zap: z}, nil
}

func (l *LoggerZap) Info(msg string, fields ...zap.Field) {
	l.zap.Info(msg, fields...)
}

func (l *LoggerZap) Error(msg string, fields ...zap.Field) {
	l.zap.Error(msg, fields...)
}

func (l *LoggerZap) Warn(msg string, fields ...zap.Field) {
	l.zap.Warn(msg, fields...)
}

func (l *LoggerZap) Debug(msg string, fields ...zap.Field) {
	l.zap.Debug(msg, fields...)
}

func (l *LoggerZap) Sync() error {
	return l.zap.Sync()
}
