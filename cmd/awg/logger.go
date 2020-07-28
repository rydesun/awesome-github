package main

import (
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/rydesun/awesome-github/awg"
	"github.com/rydesun/awesome-github/exch/config"
	"github.com/rydesun/awesome-github/exch/github"
	"github.com/rydesun/awesome-github/lib/cohttp"
	"github.com/rydesun/awesome-github/lib/errcode"
)

type LoggerConfig struct {
	Level    zapcore.Level
	Path     []string
	Encoding string
}

func setLoggers(config config.Loggers) (*zap.Logger, error) {
	defaultLoggerConfig := getLoggerConfig(config.Main)
	httpLoggerConfig := getLoggerConfig(config.Http)
	defaultLogger, err := zap.Config{
		Level:             zap.NewAtomicLevelAt(defaultLoggerConfig.Level),
		Development:       false,
		Encoding:          defaultLoggerConfig.Encoding,
		EncoderConfig:     zap.NewDevelopmentEncoderConfig(),
		OutputPaths:       defaultLoggerConfig.Path,
		ErrorOutputPaths:  defaultLoggerConfig.Path,
		DisableStacktrace: true,
	}.Build()
	if err != nil {
		return nil, err
	}
	httpLogger, err := zap.Config{
		Level:            zap.NewAtomicLevelAt(httpLoggerConfig.Level),
		Development:      false,
		Encoding:         httpLoggerConfig.Encoding,
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
		OutputPaths:      httpLoggerConfig.Path,
		ErrorOutputPaths: httpLoggerConfig.Path,
	}.Build()
	if err != nil {
		return nil, err
	}
	awg.SetDefaultLogger(defaultLogger)
	errcode.SetDefaultLogger(defaultLogger)
	github.SetDefaultLogger(defaultLogger)

	cohttp.SetDefaultLogger(httpLogger)
	return defaultLogger, nil
}

func getLoggerConfig(config config.Logger) LoggerConfig {
	level, ok := map[string]zapcore.Level{
		"debug": zap.DebugLevel,
		"info":  zap.InfoLevel,
		"warn":  zap.WarnLevel,
		"error": zap.ErrorLevel,
		"panic": zap.PanicLevel,
	}[strings.ToLower(strings.TrimSpace(config.Level))]
	if !ok {
		level = zap.InfoLevel
	}
	path := config.Path
	if len(path) == 0 {
		path = []string{"stderr"}
	}
	var encoding string
	if config.Console {
		encoding = "console"
	} else {
		encoding = "json"
	}
	return LoggerConfig{
		Level:    level,
		Path:     path,
		Encoding: encoding,
	}
}
