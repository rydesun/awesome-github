package main

import (
	"io/ioutil"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/rydesun/awesome-github/exch/config"
	"github.com/rydesun/awesome-github/exch/github"
	"github.com/rydesun/awesome-github/test/fake-github"
)

type TestEnv struct {
	testdataHolder fakeg.DataHolder
	htmlTestServer *httptest.Server
	apiTestServer  *httptest.Server
}

func (t *TestEnv) Setup() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	testdataDir := filepath.Join(wd, "../../test/testdata")
	t.testdataHolder = fakeg.NewDataHolder(testdataDir)

	htmlTestServer, err := fakeg.HtmlServer(t.testdataHolder)
	if err != nil {
		return err
	}
	apiTestServer, err := fakeg.ApiServer(t.testdataHolder)
	if err != nil {
		return err
	}
	t.htmlTestServer = htmlTestServer
	t.apiTestServer = apiTestServer
	return nil
}

func TestWorker_work(t *testing.T) {
	require := require.New(t)
	testEnv := TestEnv{}
	err := testEnv.Setup()
	require.Nil(err)

	testCases := []struct {
		config config.Config
		hasErr bool
	}{
		{
			config: config.Config{
				AccessToken: "123456",
				StartPoint: config.StartPoint{
					ID: github.RepoID{
						Owner: "tester",
						Name:  "awesome-test",
					},
				},
				Github: config.Github{
					HTMLHost: testEnv.htmlTestServer.URL,
					ApiHost:  testEnv.apiTestServer.URL,
				},
				Cli: config.Cli{
					DisableProgressBar: true,
				},
			},
			hasErr: false,
		},
		{
			config: config.Config{
				AccessToken: "123456",
				StartPoint: config.StartPoint{
					ID: github.RepoID{
						Owner: "tester",
						Name:  "awesome-test",
					},
				},
				Github: config.Github{
					HTMLHost: testEnv.htmlTestServer.URL,
					ApiHost:  testEnv.apiTestServer.URL,
				},
				Cli: config.Cli{
					// Enable progress bar
					DisableProgressBar: false,
				},
			},
			hasErr: false,
		},
		{
			config: config.Config{
				// Invalid
				AccessToken: "invalid",
				StartPoint: config.StartPoint{
					ID: github.RepoID{
						Owner: "tester",
						Name:  "awesome-test",
					},
				},
				Network: config.Net{
					RetryTime:     2,
					RetryInterval: time.Second,
				},
				Github: config.Github{
					HTMLHost: testEnv.htmlTestServer.URL,
					ApiHost:  testEnv.apiTestServer.URL,
				},
				Cli: config.Cli{
					DisableProgressBar: true,
				},
			},
			hasErr: true,
		},
		{
			config: config.Config{
				AccessToken: "123456",
				StartPoint: config.StartPoint{
					// Invalid
					ID: github.RepoID{
						Owner: "invalid",
						Name:  "invalid",
					},
				},
				Github: config.Github{
					HTMLHost: testEnv.htmlTestServer.URL,
					ApiHost:  testEnv.apiTestServer.URL,
				},
				Cli: config.Cli{
					DisableProgressBar: true,
				},
			},
			hasErr: true,
		},
		{
			config: config.Config{
				AccessToken: "123456",
				StartPoint: config.StartPoint{
					ID: github.RepoID{
						Owner: "tester",
						Name:  "awesome-test",
					},
				},
				// Invalid
				Github: config.Github{
					HTMLHost: "invalid",
					ApiHost:  "invalid",
				},
				Cli: config.Cli{
					DisableProgressBar: true,
				},
			},
			hasErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run("NONAME", func(*testing.T) {
			worker := NewWorker(ioutil.Discard, zap.NewNop())
			err = worker.Init(tc.config)
			require.Nil(err)
			err = worker.Work()
			if tc.hasErr {
				require.NotNil(err)
			} else {
				require.Nil(err)
			}
		})
	}
}
