package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func InitLogger() {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	var err error
	Log, err = config.Build()
	if err != nil {
		panic("Logger doesnt respond " + err.Error())
	}
	Log.Info("Logger started successfully")
}
