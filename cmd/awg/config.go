package main

import (
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"

	"github.com/rydesun/awesome-github/exch/config"
)

func parseConfig() (config.Config, error) {
	logger := getLogger()
	defer logger.Sync()

	// parse flags
	flagParser := config.FlagParser{}
	conf, err := flagParser.Parse()
	if err != nil {
		logger.Error("failed to read arguments", zap.Error(err))
		return conf, err
	}
	configPath := conf.ConfigPath

	// parse config file if it is specified
	if len(conf.ConfigPath) > 0 {
		const errMsg = "failed to parse config file"
		yamlParser, err := config.NewYAMLParser(conf.ConfigPath)
		if err != nil {
			return conf, err
		}
		conf, err = config.GetConfig(yamlParser)
		if err != nil {
			return conf, err
		}
		conf.ConfigPath = configPath
	}
	return conf, nil
}

func inspectConfig(conf config.Config) error {
	logger := getLogger()
	defer logger.Sync()

	confBytes, err := yaml.Marshal(config.NewProtectedConfig(conf))
	if err != nil {
		errMsg := "failed to inspect config"
		logger.Warn(errMsg, zap.Error(err))
		return err
	}
	logger.Info(string(confBytes))
	return nil
}
