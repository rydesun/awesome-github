package main

import (
	"fmt"
	"io"
	"os"

	"github.com/rydesun/awesome-github/exch/config"
)

func main() {
	// CLI
	writer := os.Stdout

	// Always parse config first.
	config, err := parseConfig(writer)
	if err != nil {
		fmt.Fprintln(writer, err)
		os.Exit(1)
	}

	// Set loggers right now.
	logger, err := setLoggers(config)
	if err != nil {
		fmt.Fprintln(writer, err)
		os.Exit(1)
	}

	// The worker takes over all writers,
	// do not write again.
	worker := NewWorker(writer, logger)
	// Must init first.
	err = worker.Init(config)
	if err != nil {
		os.Exit(1)
	}
	err = worker.Work()
	if err != nil {
		os.Exit(1)
	}
}

func parseConfig(writer io.Writer) (config.Config, error) {
	flagParser := config.FlagParser{}
	flags, err := flagParser.Parse()
	if err != nil {
		err = fmt.Errorf("Failed to parse cmd flags. %v", err)
		return config.Config{}, err
	}
	configPath := flags.ConfigPath

	yamlParser, err := config.NewYAMLParser(flags.ConfigPath)
	if err != nil {
		err = fmt.Errorf("Failed to parse config files. %v", err)
		return config.Config{}, err
	}

	conf, err := config.GetConfig(yamlParser)
	if err != nil {
		err = fmt.Errorf("Failed to parse config files. %v", err)
		return config.Config{}, err
	}
	flags.ConfigPath = configPath
	return conf, nil
}
