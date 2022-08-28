// Copyright (c) 554949297@qq.com . 2022-2022 . All rights reserved

package g3

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"path"
	"time"
)

const (
	DefaultLogDir = "logs"
)

func init() {
}

func DefaultLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + l.String() + "]")
}

func DefaultTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func NewLogger(name string, caller bool) *zap.Logger {
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   path.Join(HomeDir(), DefaultLogDir, name+".log"),
		MaxSize:    20, // megabytes
		MaxBackups: 10,
		MaxAge:     30, // days
	})

	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = DefaultTimeEncoder
	cfg.EncodeLevel = DefaultLevelEncoder

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(cfg),
		w,
		zap.InfoLevel,
	)

	opts := make([]zap.Option, 0)
	opts = append(opts, zap.AddStacktrace(zap.ErrorLevel))
	if caller {
		opts = append(opts, zap.AddCaller())
	}

	return zap.New(core, opts...)
}
