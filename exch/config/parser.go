package config

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sys/unix"

	"github.com/rydesun/awesome-github/lib/errcode"
)

type ConfigParser interface {
	Parse() (Config, error)
}

func GetConfig(cp ConfigParser) (Config, error) {
	config, err := cp.Parse()
	if err != nil {
		return Config{}, err
	}
	accessToken := os.Getenv("GITHUB_ACCESS_TOKEN")
	if len(accessToken) > 0 {
		config.AccessToken = accessToken
	}
	config.Output.Path, err = filepath.Abs(config.Output.Path)
	if err != nil {
		return Config{}, err
	}
	if !CheckFileWritable(config.Output.Path) {
		errMsg := fmt.Sprintf("Failed to Write file. Invalid output path: %s", config.Output.Path)
		return Config{}, errcode.New(errMsg, ErrCodeParameter, ErrScope, nil)
	}
	for i, p := range config.Log.Main.Path {
		p, err = filepath.Abs(p)
		if !CheckFileWritable(p) {
			errMsg := fmt.Sprintf("Failed to Write file. Invalid output path: %s", p)
			return Config{}, errcode.New(errMsg,
				ErrCodeParameter, ErrScope, nil)
		}
		config.Log.Main.Path[i] = p
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

// Check this file path is writable.
// Will create directories first.
func CheckFileWritable(path string) bool {
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return false
	}
	return unix.Access(dir, unix.W_OK) == nil
}
