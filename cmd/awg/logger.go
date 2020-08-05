package main

import (
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/rydesun/awesome-github/awg"
	"github.com/rydesun/awesome-github/exch/config"
	"github.com/rydesun/awesome-github/exch/github"
	"github.com/rydesun/awesome-github/lib/cohttp"
	"github.com/rydesun/awesome-github/lib/errcode"
)

type LoggerConfig struct {
	Level    zapcore.Level
	Path     string
	Encoding string
}

func setLoggers(config config.Loggers) (*zap.Logger, error) {
	defaultLoggerConfig := getLoggerConfig(config.Main)
	httpLoggerConfig := getLoggerConfig(config.Http)

	enc := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      "C",
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})

	defaultWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   defaultLoggerConfig.Path,
		MaxSize:    20,
		MaxBackups: 1,
		MaxAge:     28,
	})
	defaultLogger := zap.New(
		zapcore.NewCore(enc, defaultWriter, zap.NewAtomicLevelAt(
			defaultLoggerConfig.Level)), zap.AddCaller())

	var httpLogger *zap.Logger
	if httpLoggerConfig.Path == defaultLoggerConfig.Path {
		httpLogger = zap.New(
			zapcore.NewCore(enc, defaultWriter, zap.NewAtomicLevelAt(
				httpLoggerConfig.Level)), zap.AddCaller())
	} else {
		httpWriter := zapcore.AddSync(&lumberjack.Logger{
			Filename:   httpLoggerConfig.Path,
			MaxSize:    5,
			MaxBackups: 1,
			MaxAge:     28,
		})
		httpLogger = zap.New(
			zapcore.NewCore(enc, httpWriter, zap.NewAtomicLevelAt(
				httpLoggerConfig.Level)), zap.AddCaller())
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
	return LoggerConfig{
		Level: level,
		Path:  config.Path,
	}
}
