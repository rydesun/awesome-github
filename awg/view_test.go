package awg

import (
	"encoding/json"
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/rydesun/awesome-github/exch/github"
	"github.com/rydesun/awesome-github/test/fake-github"
)

func TestRepo_Aggregate(t *testing.T) {
	require := require.New(t)
	wd, err := os.Getwd()
	require.Nil(err)
	testdataDir := path.Join(wd, "../test/testdata")
	testdataHolder := fakeg.NewDataHolder(testdataDir)
	raw, err := testdataHolder.GetJsonRepo("goxjs", "glfw")
	require.Nil(err)
	gbRepo := &github.Repo{}
	err = json.Unmarshal(raw, gbRepo)
	require.Nil(err)

	repo := Repo{}
	err = repo.Aggregate(gbRepo)
	require.Nil(err)

	lastCommitTime, err := time.Parse(time.RFC3339, "2019-07-15T19:40:41Z")
	require.Nil(err)
	require.Equal(Repo{
		LastCommit:  lastCommitTime,
		Star:        65,
		Description: "Go cross-platform glfw library for creating an OpenGL context and receiving events.",
	}, repo)
}
