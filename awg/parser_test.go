package awg

import (
	"os"
	"path"
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
	testdataDir := path.Join(wd, "../test/testdata")
	testdataHolder := fakeg.NewDataHolder(testdataDir)
	htmlTestServer, err := fakeg.ApiServer(testdataHolder)
	require.Nil(err)
	apiTestServer, err := fakeg.ApiServer(testdataHolder)
	require.Nil(err)
	htmlReadme, err := testdataHolder.GetHtmlAwesomeReadme()
	require.Nil(err)
	htmlClient := cohttp.NewClient(*htmlTestServer.Client(), 16, 20, nil)
	apiClient := cohttp.NewClient(*apiTestServer.Client(), 16, 20, nil)
	gbClient, err := github.NewClient(htmlClient, apiClient,
		github.ClientOption{
			HTMLHost:    htmlTestServer.URL,
			HTMLPathPre: github.HTMLPathPre,
			APIHost:     apiTestServer.URL,
			ApiPathPre:  github.APIPathPre,
		})
	require.Nil(err)
	client, err := NewClient(gbClient)
	require.Nil(err)
	t.Run("", func(t *testing.T) {
		reporter := &Reporter{}
		awesomeParser, err := NewParser(string(htmlReadme), client, reporter)
		require.Nil(err)
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
						Star:        319,
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
						Star:        1,
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
						Star:        65,
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
	// Test invalid README.md
	t.Run("invalid-1", func(t *testing.T) {
		reporter := &Reporter{}
		awesomeParser, err := NewParser("", client, reporter)
		require.Equal(nil, err)
		_, err = awesomeParser.Gather()
		require.NotEqual(nil, err)
	})
	// Test invalid README.md
	t.Run("invalid-2", func(t *testing.T) {
		reporter := &Reporter{}
		awesomeParser, err := NewParser("<h2><li>", client, reporter)
		require.Equal(nil, err)
		_, err = awesomeParser.Gather()
		require.NotEqual(nil, err)
	})
}
