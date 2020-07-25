package config

import "os"

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
	err = config.Validate()
	return config, err
}
