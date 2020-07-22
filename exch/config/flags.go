package config

import (
	"flag"

	"github.com/rydesun/awesome-github/exch/github"
)

type FlagParser struct {
}

func (p *FlagParser) Parse() (Config, error) {
	configPath := flag.String("config", "", "yaml config path")
	user := flag.String("user", "avelino", "awesome repository owner")
	name := flag.String("name", "awesome-go", "awesome repository name")
	maxConcurrent := flag.Int("max-concurrent", 6, "max concurrent requests")
	logRespHead := flag.Int("log-head", 200, "truncated response in log")
	accessToken := flag.String("token", "", "your github access token")
	flag.Parse()
	return Config{
		ConfigPath:  *configPath,
		AccessToken: *accessToken,
		StartPoint: StartPoint{
			Id: github.RepoID{
				Owner: *user,
				Name:  *name,
			},
		},
		MaxConcurrent: *maxConcurrent,
		LogRespHead:   *logRespHead,
	}, nil
}
