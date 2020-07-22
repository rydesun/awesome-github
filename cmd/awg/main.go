package main

import (
	"log"

	"go.uber.org/zap"
)

func main() {
	conf, err := parseConfig()
	if err != nil {
		log.Panic(err)
	}
	err = inspectConfig(conf)
	if err != nil {
		log.Panic(err)
	}

	setLoggers(conf)
	logger := getLogger()

	reporter := newReporter()
	client, err := newClient(conf, reporter)
	if err != nil {
		logger.Panic("failed to create github client", zap.Error(err))
	}

	err = work(client, conf, reporter)
	if err != nil {
		logger.Panic("failed to analysize awesome repositories", zap.Error(err))
	}
}
