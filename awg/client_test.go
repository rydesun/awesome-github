package awg

import (
	"context"
	"encoding/json"
	"flag"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/rydesun/awesome-github/exch/github"
	"github.com/rydesun/awesome-github/lib/cohttp"
	"github.com/rydesun/awesome-github/test/fake-github"
)

var realSrc bool
var accessToken string

func init() {
	flag.BoolVar(&realSrc, "real", false, "fetch data from real github")
	flag.StringVar(&accessToken, "token", "", "your github access token")
}

type ClientTestEnv struct {
	awgClient      *Client
	testdataHolder fakeg.DataHolder
	testServer     *httptest.Server
}

func (t *ClientTestEnv) Setup() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	testdataDir := filepath.Join(wd, "../test/testdata")
	testdataHolder := fakeg.NewDataHolder(testdataDir)
	var gbClient *github.Client
	if !realSrc {
		testServer, err := fakeg.ApiServer(testdataHolder)
		if err != nil {
			return err
		}
		t.testServer = testServer
		gbClient, err = github.NewClient(
			nil,
			cohttp.NewClient(*testServer.Client(), 16, 0, time.Second, 20, nil),
			github.ClientOption{
				APIHost:     testServer.URL,
				ApiPathPre:  github.APIPathPre,
				AccessToken: "123456",
			})
		if err != nil {
			return err
		}
	} else {
		gbClient, err = github.NewClient(
			nil,
			cohttp.NewClient(http.Client{}, 16, 0, time.Second, 20, nil),
			github.ClientOption{
				APIHost:     github.APIHost,
				ApiPathPre:  github.APIPathPre,
				AccessToken: accessToken,
			})
	}
	if err != nil {
		return err
	}
	awgClient, err := NewClient(gbClient)
	if err != nil {
		return err
	}

	t.testdataHolder = testdataHolder
	t.awgClient = awgClient
	return nil
}

func TestClient_GetUser(t *testing.T) {
	require := require.New(t)
	testEnv := ClientTestEnv{}
	err := testEnv.Setup()
	require.Nil(err)

	if !realSrc {
		user, err := testEnv.awgClient.GetUser()
		require.Nil(err)
		require.Equal("tester", user.Name)
		require.Equal(5000, user.RateLimit.Total)
		require.Equal(4999, user.RateLimit.Remaining)
		require.False(user.RateLimit.ResetAt.IsZero())

		testEnv.testServer.Close()
		_, err = testEnv.awgClient.GetUser()
		require.NotNil(err)
	} else {
		user, err := testEnv.awgClient.GetUser()
		require.Nil(err)
		require.NotNil(user.Name)
		require.Greater(0, user.RateLimit.Total)
		require.False(user.RateLimit.ResetAt.IsZero())
	}
}

func TestClient_Fill(t *testing.T) {
	require := require.New(t)
	testEnv := ClientTestEnv{}
	err := testEnv.Setup()
	require.Nil(err)

	testCases := []struct {
		user   string
		name   string
		hasErr bool
	}{
		{
			user: "antchfx",
			name: "xpath",
		},
		{
			user:   "invalidUser",
			name:   "invalidName",
			hasErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.user+"/"+tc.name, func(t *testing.T) {
			awesomeRepo := AwesomeRepo{
				Repo: Repo{
					ID: github.RepoID{
						Owner: tc.user,
						Name:  tc.name,
					},
					Owner:       tc.user,
					AwesomeName: tc.name,
				},
			}
			err = testEnv.awgClient.Fill(context.Background(), &awesomeRepo)
			if err != nil {
				if tc.hasErr {
					// expected error
					return
				}
				t.Fatal(err)
			}
			if !realSrc {
				content, err := testEnv.testdataHolder.GetJsonRepo(tc.user, tc.name)
				require.Nil(err)
				expectedRepo := github.Repo{}
				_ = json.Unmarshal(content, &expectedRepo)
				require.Equal(expectedRepo.Data.Repository.Stargazers.TotalCount,
					awesomeRepo.Star)
				require.Equal(expectedRepo.Data.Repository.DefaultBranchRef.Target.History.Edges[0].Node.CommittedDate,
					awesomeRepo.LastCommit)
				require.Equal(expectedRepo.Data.Repository.Description,
					awesomeRepo.Description)
			} else {
				require.Less(0, awesomeRepo.Star)
				require.NotEmpty(awesomeRepo.LastCommit)
				require.NotEmpty(awesomeRepo.Description)
			}
		})
	}
}
