package config

import (
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/rydesun/awesome-github/lib/errcode"
)

type YAMLParser struct {
	fpath string
}

func NewYAMLParser(fpath string) (*YAMLParser, error) {
	_, err := os.Stat(fpath)
	if err != nil {
		return nil, err
	}
	return &YAMLParser{
		fpath: fpath,
	}, nil
}

func (p *YAMLParser) Parse() (Config, error) {
	config := Config{}
	raw, err := ioutil.ReadFile(p.fpath)
	if err != nil {
		return config, err
	}
	err = yaml.UnmarshalStrict(raw, &config)
	if err != nil {
		return config, err
	}
	config.ConfigPath = p.fpath
	path := config.StartPoint.Path
	sliceStr := strings.Split(path, "/")
	if len(sliceStr) != 2 {
		errcode.New("Invaild path",
			ErrCodeParameter, ErrScope, []string{"path"})
	}
	config.StartPoint.Id.Owner = sliceStr[0]
	config.StartPoint.Id.Name = sliceStr[1]
	return config, nil
}
