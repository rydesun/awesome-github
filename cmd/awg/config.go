package main

import (
	"fmt"
	"os"

	"go.uber.org/zap"

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
	writer := os.Stdout
	fmt.Fprintf(writer, "config file: %s\n", conf.ConfigPath)
	fmt.Fprintf(writer, "main log files: %s\n", conf.Log.Main.Path)
	fmt.Fprintf(writer, "http log files: %s\n", conf.Log.Http.Path)
	fmt.Fprintf(writer, "output file: %s\n", conf.Output.Path)
	return nil
}
