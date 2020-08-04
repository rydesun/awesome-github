package config

import (
	"strings"
	"time"

	"github.com/go-playground/validator/v10"

	"github.com/rydesun/awesome-github/exch/github"
	"github.com/rydesun/awesome-github/lib/errcode"
)

type Config struct {
	ConfigPath    string `yaml:"config"`
	AccessToken   string `yaml:"access_token" validate:"required"`
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
	Path          string `yaml:"path" validate:"required"`
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
	Level   string   `yaml:"level"`
	Path    []string `yaml:"path"`
	Console bool     `yaml:"console"`
}

func (c *Config) Validate() error {
	validate := validator.New()
	err := validate.Struct(c)
	if err != nil {
		return err
	}
	if len(c.Github.HTMLHost) > 0 {
		err := validate.Var(c.Github.HTMLHost, "url")
		if err != nil {
			errMsg := "Invalid github html host"
			return errcode.New(errMsg, ErrCodeParameter, ErrScope,
				[]string{"githubHTMLHost"})
		}
	}
	if len(c.Github.ApiHost) > 0 {
		err := validate.Var(c.Github.ApiHost, "url")
		if err != nil {
			errMsg := "Invalid github api host"
			return errcode.New(errMsg, ErrCodeParameter, ErrScope,
				[]string{"githubAPIHost"})
		}
	}
	return nil
}

func SplitID(id string) (owner, name string, err error) {
	sliceStr := strings.Split(id, "/")
	if len(sliceStr) != 2 {
		return "", "", errcode.New("Invaild path",
			ErrCodeParameter, ErrScope, []string{"path"})
	}
	return sliceStr[0], sliceStr[1], nil
}
