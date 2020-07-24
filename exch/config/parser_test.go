package config

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/rydesun/awesome-github/exch/github"
	"github.com/stretchr/testify/require"
)

const testAccessToken = "123456"

func TestMain(m *testing.M) {
	oldAccessToken := os.Getenv("GITHUB_ACCESS_TOKEN")
	err := os.Setenv("GITHUB_ACCESS_TOKEN", testAccessToken)
	if err != nil {
		log.Fatalln(err)
	}
	code := m.Run()
	os.Setenv("GITHUB_ACCESS_TOKEN", oldAccessToken)
	os.Exit(code)
}

func TestGetConfig(t *testing.T) {
	require := require.New(t)

	path := "../../configs/config.yaml"
	yamlParser, err := NewYAMLParser(path)
	require.Nil(err)
	actual, err := GetConfig(yamlParser)
	require.Nil(err)

	expected := Config{
		ConfigPath:    path,
		AccessToken:   testAccessToken,
		MaxConcurrent: 3,
		StartPoint: StartPoint{
			Path: "avelino/awesome-go",
			ID: github.RepoID{
				Owner: "avelino",
				Name:  "awesome-go",
			},
		},
		Network: Net{
			RetryTime:     2,
			RetryInterval: time.Second,
		},
		Output: Output{
			Path: "./awg.json",
		},
		Log: Loggers{
			Main: Logger{
				Path: []string{"/tmp/awesome-github.log"},
			},
			Http: Logger{
				Path: []string{"/tmp/awesome-github.log"},
			},
		},
	}
	require.Equal(expected, actual)

	// Invalid config path
	_, err = NewYAMLParser("")
	require.NotNil(err)
	_, err = NewYAMLParser("a/b/c/d")
	require.NotNil(err)
}
