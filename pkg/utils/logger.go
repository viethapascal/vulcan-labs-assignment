package utils

import (
	"go.uber.org/zap"
)

var appLogger map[string]*zap.SugaredLogger = make(map[string]*zap.SugaredLogger)
var BaseLogger *zap.SugaredLogger

func init() {
	BaseLogger = Logger("main-app")
}

func Logger(name string) *zap.SugaredLogger {
	if logger, ok := appLogger[name]; ok {
		return logger
	}
	l, _ := zap.NewDevelopment()

	appLogger[name] = l.Named(name).Sugar()
	//logger, _ := zap.NewDevelopment()
	//sugar := logger.Sugar()
	return appLogger[name]
}
