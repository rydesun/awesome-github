package config

import (
	"time"

	"github.com/rydesun/awesome-github/exch/github"
)

type Config struct {
	ConfigPath    string `yaml:"config"`
	AccessToken   string `yaml:"access_token"`
	MaxConcurrent int    `yaml:"max_concurrent"`
	LogRespHead   int    `yaml:"log_resp_head"`
	StartPoint    `yaml:"start_point"`
	Network       Net     `yaml:"network"`
	Output        Output  `yaml:"output"`
	Github        Github  `yaml:"github"`
	Cli           Cli     `yaml:"cli"`
	Log           Loggers `yaml:"log"`
}

type StartPoint struct {
	Path          string
	ID            github.RepoID
	SectionFilter []string
}

type Net struct {
	RetryTime     int           `yaml:"retry_time"`
	RetryInterval time.Duration `yaml:"retry_interval"`
}

type Output struct {
	Path string `yaml:"path"`
}

func NewProtectedConfig(config Config) Config {
	config.AccessToken = "<GitHub Personal Access Token>"
	return config
}

type Github struct {
	HTMLHost string `yaml:"html_host"`
	ApiHost  string `yaml:"api_host"`
}

type Cli struct {
	DisableProgressBar bool `yaml:"disable_progress_bar"`
}

type Loggers struct {
	Main Logger `yaml:"main"`
	Http Logger `yaml:"http"`
}

type Logger struct {
	Level   string
	Path    []string
	Console bool
}
