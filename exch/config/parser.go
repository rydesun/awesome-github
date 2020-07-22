package config

import "os"

type ConfigParser interface {
	Parse() (Config, error)
}

func GetConfig(cp ConfigParser) (config Config, err error) {
	config, err = cp.Parse()
	if err != nil {
		return
	}
	accessToken := os.Getenv("GITHUB_ACCESS_TOKEN")
	if len(accessToken) > 0 {
		config.AccessToken = accessToken
	}
	return
}
