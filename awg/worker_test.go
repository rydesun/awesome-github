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

func TestWorkflow(t *testing.T) {
	require := require.New(t)
	wd, err := os.Getwd()
	require.Nil(err)
	testdataDir := path.Join(wd, "../test/testdata")
	testdataHolder := fakeg.NewDataHolder(testdataDir)
	htmlServer, err := fakeg.HtmlServer(testdataHolder)
	require.Nil(err)
	apiServer, err := fakeg.ApiServer(testdataHolder)
	require.Nil(err)
	gbClient, err := github.NewClient(
		cohttp.NewClient(*htmlServer.Client(), 16, 2, time.Second, 20, nil),
		cohttp.NewClient(*apiServer.Client(), 16, 2, time.Second, 20, nil),
		github.ClientOption{
			HTMLHost:    htmlServer.URL,
			HTMLPathPre: github.HTMLPathPre,
			APIHost:     apiServer.URL,
			ApiPathPre:  github.APIPathPre,
		})
	require.Nil(err)
	client, err := NewClient(gbClient)
	require.Nil(err)
	result, err := Workflow(client, nil, github.RepoID{Owner: "tester", Name: "awesome-test"},
		RateLimit{
			Total:     100000,
			Remaining: 100000,
		})
	require.Nil(err)
	require.Less(0, len(result))

	// Test invalid, should have a error.
	result, err = Workflow(client, nil, github.RepoID{Owner: "invalid", Name: "invalid"},
		RateLimit{
			Total:     100000,
			Remaining: 100000,
		})
	require.NotNil(err)
}
