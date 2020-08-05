package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/rydesun/awesome-github/exch/config"
	"github.com/rydesun/awesome-github/web/app"
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
			{
				Name:   "view",
				Usage:  "View awesome readme in browser",
				Action: view,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:      "script",
						Usage:     "Embedded script",
						Required:  true,
						TakesFile: true,
					},
					&cli.StringFlag{
						Name:  "listen",
						Usage: "Listen address",
						Value: "127.0.0.1:3000",
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

func view(c *cli.Context) error {
	// CLI
	writer := os.Stdout
	depressLoggers()

	if c.Args().Len() == 0 {
		fmt.Fprintln(writer, "Awesome list name is missing")
		cli.OsExiter(1)
	}

	datafile := c.Args().Get(0)
	data, err := LoadOutputFile(datafile)
	if err != nil {
		fmt.Fprintln(writer, strerr(err))
		cli.OsExiter(1)
	}
	if !data.IsValid() {
		fmt.Fprintln(writer, "Invalid data")
		cli.OsExiter(1)
	}

	router, err := app.NewRouter(c.String("listen"))
	if err != nil {
		fmt.Fprintln(writer, strerr(err))
		cli.OsExiter(1)
	}
	scriptPath := c.String("script")
	dataPath := datafile
	fmt.Fprintln(writer, "[1/2] Fetching remote readme page...")
	err = router.Init(data.AwesomeList, scriptPath, dataPath)
	if err != nil {
		fmt.Fprintln(writer, strerr(err))
		cli.OsExiter(1)
	}
	fmt.Fprintf(writer, "[2/2] Serve at http://%s\n", c.String("listen"))
	err = router.Route()
	if err != nil {
		fmt.Fprintln(writer, strerr(err))
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
