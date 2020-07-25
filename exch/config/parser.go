package config

import (
	"os"
	"path/filepath"
)

type ConfigParser interface {
	Parse() (Config, error)
}

func GetConfig(cp ConfigParser) (Config, error) {
	config, err := cp.Parse()
	if err != nil {
		return Config{}, nil
	}
	accessToken := os.Getenv("GITHUB_ACCESS_TOKEN")
	if len(accessToken) > 0 {
		config.AccessToken = accessToken
	}
	config.Output.Path, err = filepath.Abs(config.Output.Path)
	if err != nil {
		return Config{}, err
	}
	for i, p := range config.Log.Main.Path {
		config.Log.Main.Path[i], err = filepath.Abs(p)
		if err != nil {
			return Config{}, err
		}
	}
	for i, p := range config.Log.Http.Path {
		config.Log.Http.Path[i], err = filepath.Abs(p)
		if err != nil {
			return Config{}, err
		}
	}
	err = config.Validate()
	return config, err
}
