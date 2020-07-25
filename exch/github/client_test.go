package github

import (
	"net/http/httptest"
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/rydesun/awesome-github/lib/cohttp"
	"github.com/rydesun/awesome-github/test/fake-github"
)

type ClientTestEnv struct {
	apiTestServer  *httptest.Server
	htmlTestServer *httptest.Server
	testdataHolder fakeg.DataHolder
}

func (t *ClientTestEnv) Setup() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	testdataDir := path.Join(wd, "../../test/testdata")
	t.testdataHolder = fakeg.NewDataHolder(testdataDir)

	t.apiTestServer, err = fakeg.ApiServer(t.testdataHolder)
	t.htmlTestServer, err = fakeg.HtmlServer(t.testdataHolder)
	return err
}

func TestClient_GetUser(t *testing.T) {
	require := require.New(t)
	testEnv := ClientTestEnv{}
	err := testEnv.Setup()
	require.Nil(err)

	testServer := testEnv.apiTestServer
	apiClient := cohttp.NewClient(*testServer.Client(), 16, 0, time.Second, 20, nil)

	testCases := []struct {
		token  string
		hasErr bool
	}{
		{
			token:  "123456",
			hasErr: false,
		},
		{
			token:  "invalid",
			hasErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.token, func(t *testing.T) {
			client, err := NewClient(nil, apiClient,
				ClientOption{
					APIHost:     testServer.URL,
					ApiPathPre:  APIPathPre,
					AccessToken: tc.token,
				})
			require.Nil(err)
			user, err := client.GetUser()
			if tc.hasErr {
				require.NotNil(err)
			} else {
				require.Nil(err)
				require.NotNil(user)
			}
		})
	}
}

func TestClient_GetRepo(t *testing.T) {
	require := require.New(t)
	testEnv := ClientTestEnv{}
	err := testEnv.Setup()
	require.Nil(err)

	testServer := testEnv.apiTestServer
	apiClient := cohttp.NewClient(*testServer.Client(), 16, 0, time.Second, 20, nil)

	testCases := []struct {
		repoID RepoID
		hasErr bool
	}{
		{
			repoID: RepoID{
				Owner: "antchfx",
				Name:  "xpath",
			},
			hasErr: false,
		},
		{
			repoID: RepoID{
				Owner: "invalid",
				Name:  "invalid",
			},
			hasErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.repoID.String(), func(t *testing.T) {
			client, err := NewClient(nil, apiClient,
				ClientOption{
					APIHost:    testServer.URL,
					ApiPathPre: APIPathPre,
				})
			require.Nil(err)
			user, err := client.GetRepo(tc.repoID)
			if tc.hasErr {
				require.NotNil(err)
			} else {
				require.Nil(err)
				require.NotNil(user)
			}
		})
	}
}

func TestClient_GetHTMLReadme(t *testing.T) {
	require := require.New(t)
	testEnv := ClientTestEnv{}
	err := testEnv.Setup()
	require.Nil(err)

	testServer := testEnv.htmlTestServer
	htmlClient := cohttp.NewClient(*testServer.Client(), 16, 0, time.Second, 20, nil)

	testCases := []struct {
		repoID RepoID
		hasErr bool
	}{
		{
			repoID: RepoID{
				Owner: "tester",
				Name:  "awesome-test",
			},
			hasErr: false,
		},
		{
			repoID: RepoID{
				Owner: "invalid",
				Name:  "invalid",
			},
			hasErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.repoID.String(), func(t *testing.T) {
			client, err := NewClient(htmlClient, nil,
				ClientOption{
					HTMLHost:    testServer.URL,
					HTMLPathPre: HTMLPathPre,
				})
			require.Nil(err)
			user, err := client.GetHTMLReadme(tc.repoID)
			if tc.hasErr {
				require.NotNil(err)
			} else {
				require.Nil(err)
				require.NotNil(user)
			}
		})
	}
}
