package config

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/rydesun/awesome-github/lib/errcode"
)

type YAMLParser struct {
	fpath string
}

func NewYAMLParser(fpath string) (*YAMLParser, error) {
	if len(fpath) == 0 {
		return nil, errcode.New("Invalid config path",
			ErrCodeParameter, ErrScope, nil)
	}
	fpath, err := filepath.Abs(fpath)
	if err != nil {
		return nil, err
	}
	_, err = os.Stat(fpath)
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
	owner, name, err := SplitID(config.StartPoint.Path)
	if err != nil {
		return config, err
	}
	config.StartPoint.ID.Owner = owner
	config.StartPoint.ID.Name = name
	return config, nil
}
