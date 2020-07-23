package main

import (
	"path/filepath"
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

func setLoggers(config config.Config) (err error) {
	configPath := config.ConfigPath
	defaultLoggerConfig := getLoggerConfig(configPath, config.Log.Main)
	httpLoggerConfig := getLoggerConfig(configPath, config.Log.Http)
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
		return
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
		return
	}
	setDefaultLogger(defaultLogger)
	awg.SetDefaultLogger(defaultLogger)
	errcode.SetDefaultLogger(defaultLogger)
	github.SetDefaultLogger(defaultLogger)

	cohttp.SetDefaultLogger(httpLogger)
	return
}

func getLoggerConfig(configPath string, config config.Logger) LoggerConfig {
	configDir := filepath.Dir(configPath)
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
	for i, p := range path {
		p = filepath.Clean(p)
		if !strings.HasPrefix(p, "/") {
			p = filepath.Join(configDir, p)
		}
		path[i] = p
	}
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
