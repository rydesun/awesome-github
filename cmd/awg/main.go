package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/rydesun/awesome-github/exch/config"
)

func main() {
	app := &cli.App{
		Name:  "awg",
		Usage: "Awesome GitHub repositories",
		Commands: []*cli.Command{
			{
				Name:   "fetch",
				Usage:  "Fetch data about awesome repositories",
				Action: fetch,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:      "config",
						Aliases:   []string{"c"},
						Usage:     "YAML config",
						Required:  true,
						TakesFile: true,
					},
				},
			},
		},
		EnableBashCompletion: true,
	}
	app.Run(os.Args)
}

func fetch(c *cli.Context) error {
	// CLI
	writer := os.Stdout
	configPath := c.String("config")

	// Always parse config first.
	config, err := parseConfig(configPath)
	if err != nil {
		fmt.Fprintln(writer, err)
		cli.OsExiter(1)
	}

	// Set loggers right now.
	logger, err := setLoggers(config.Log)
	if err != nil {
		fmt.Fprintln(writer, err)
		cli.OsExiter(1)
	}

	// The worker takes over all writers,
	// do not write again.
	worker := NewWorker(writer, logger)
	// Must init first.
	err = worker.Init(config)
	if err != nil {
		cli.OsExiter(1)
	}
	err = worker.Work()
	if err != nil {
		cli.OsExiter(1)
	}
	return nil
}

func parseConfig(configPath string) (config.Config, error) {
	yamlParser, err := config.NewYAMLParser(configPath)
	if err != nil {
		err = fmt.Errorf("Failed to parse config files. %v", err)
		return config.Config{}, err
	}

	conf, err := config.GetConfig(yamlParser)
	if err != nil {
		err = fmt.Errorf("Failed to parse config files. %v", strerr(err))
		return config.Config{}, err
	}
	return conf, nil
}
