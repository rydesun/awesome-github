package awg

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/rydesun/awesome-github/exch/github"
	"github.com/rydesun/awesome-github/lib/cohttp"
	"github.com/rydesun/awesome-github/test/fake-github"
)

func TestParser_Gather(t *testing.T) {
	require := require.New(t)
	wd, err := os.Getwd()
	require.Nil(err)
	testdataDir := filepath.Join(wd, "../test/testdata")
	testdataHolder := fakeg.NewDataHolder(testdataDir)
	apiTestServer, err := fakeg.ApiServer(testdataHolder)
	require.Nil(err)
	htmlReadme, err := testdataHolder.GetHtmlAwesomeReadme()
	require.Nil(err)
	apiClient := cohttp.NewClient(*apiTestServer.Client(), 16, 0, time.Second, 20, nil)
	gbClient, err := github.NewClient(nil, apiClient,
		github.ClientOption{
			APIHost:    apiTestServer.URL,
			ApiPathPre: github.APIPathPre,
		})
	require.Nil(err)
	client, err := NewClient(gbClient)
	require.Nil(err)
	t.Run("", func(t *testing.T) {
		reporter := &Reporter{}
		awesomeParser := NewParser(string(htmlReadme), client, reporter, RateLimit{
			Total:     100000,
			Remaining: 100000,
		})
		awesomeRepos, err := awesomeParser.Gather()
		require.Nil(err)

		lastCommitTime, err := time.Parse(time.RFC3339, "2019-07-15T19:40:41Z")
		require.Nil(err)
		expect := map[string][]*AwesomeRepo{
			"XML": {
				&AwesomeRepo{
					Repo: Repo{
						ID: github.RepoID{
							Owner: "antchfx",
							Name:  "xpath",
						},
						Owner:       "antchfx",
						AwesomeName: "xpath",
						Link:        "https://github.com/antchfx/xpath",
						Watch:       8,
						Star:        319,
						Fork:        36,
						LastCommit:  lastCommitTime,
						Description: "XPath package for Golang, supports HTML, XML, JSON document query.",
					},
					AwesomeDesc: "xpath - XPath package for Go.",
				},
			},
			"OpenGL": {
				{
					Repo: Repo{
						ID: github.RepoID{
							Owner: "technohippy",
							Name:  "go-glmatrix",
						},
						Owner:       "technohippy",
						AwesomeName: "go-glmatrix",
						Link:        "https://github.com/technohippy/go-glmatrix",
						Watch:       1,
						Star:        1,
						Fork:        0,
						LastCommit:  lastCommitTime,
						Description: "go-glmatrix is a golang version of glMatrix, which is \"designed to perform vector and matrix operations stupidly fast\".",
					},
					AwesomeDesc: `go-glmatrix - Go port of glMatrix library.`,
				},
				{
					Repo: Repo{
						ID: github.RepoID{
							Owner: "goxjs",
							Name:  "glfw",
						},
						Owner:       "goxjs",
						AwesomeName: "goxjs/glfw",
						Link:        "https://github.com/goxjs/glfw",
						Watch:       6,
						Star:        65,
						Fork:        14,
						LastCommit:  lastCommitTime,
						Description: "Go cross-platform glfw library for creating an OpenGL context and receiving events.",
					},
					AwesomeDesc: "goxjs/glfw - Go cross-platform glfw library for creating an OpenGL context and receiving events.",
				},
			},
		}
		require.Equal(expect, awesomeRepos)
		require.Equal([]github.RepoID{
			{
				Owner: "randominvaliduser",
				Name:  "repo",
			}}, reporter.GetInvalidRepo())
	})

	testcases := []struct {
		invalidReadme string
	}{
		{invalidReadme: ""},
		{invalidReadme: "<h2><li>"},
		{invalidReadme: `<html>
		  <div id="readme">
		  <article class="markdown-body">
		    <h2>invalid</h2>
		    <li>invalid</li>`},
	}
	for i, tc := range testcases {
		t.Run("invalid_readme_"+strconv.Itoa(i), func(t *testing.T) {
			reporter := &Reporter{}
			awesomeParser := NewParser(tc.invalidReadme, client, reporter, RateLimit{
				Total:     100000,
				Remaining: 100000,
			})
			_, err = awesomeParser.Gather()
			require.NotEqual(nil, err)
		})
	}
	// Test ratelimit
	t.Run("invalid_ratelimit", func(t *testing.T) {
		reporter := &Reporter{}
		awesomeParser := NewParser(string(htmlReadme), client, reporter, RateLimit{
			Total:     100000,
			Remaining: 0,
		})
		_, err = awesomeParser.Gather()
		require.NotNil(err)
	})
	// Test invalid network
	invalidGbClient, _ := github.NewClient(nil, apiClient,
		github.ClientOption{
			// Invalid network for GitHub API
			APIHost:    "https://127.127.127.127:12345",
			ApiPathPre: github.APIPathPre,
		})
	invalidClient, _ := NewClient(invalidGbClient)
	t.Run("invalid_network", func(t *testing.T) {
		reporter := &Reporter{}
		awesomeParser := NewParser(string(htmlReadme), invalidClient, reporter, RateLimit{
			Total:     100000,
			Remaining: 100000,
		})
		_, err = awesomeParser.Gather()
		require.NotNil(err)
		require.True(cohttp.IsNetowrkError(err))
	})
}
