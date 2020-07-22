package main

import (
	"sync"

	"go.uber.org/zap"
)

var defaultLogger *zap.Logger

func setDefaultLogger(logger *zap.Logger) {
	defaultLogger = logger
}

func getLogger() *zap.Logger {
	if defaultLogger != nil {
		return defaultLogger
	}
	var once sync.Once
	once.Do(func() {
		defaultLogger, _ = zap.NewProduction()
	})
	return defaultLogger
}
