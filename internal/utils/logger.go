package utils

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func InitLogger() {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(config)
	consoleEncoder := zapcore.NewConsoleEncoder(config)

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, zapcore.AddSync(os.Stdout), zap.InfoLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zap.InfoLevel),
	)

	Log = zap.New(core, zap.AddCaller())
}
