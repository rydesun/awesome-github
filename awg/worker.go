package awg

import (
	"github.com/rydesun/awesome-github/exch/github"
)

func Workflow(client *Client, reporter *Reporter, awesomeID github.RepoID,
	ratelimit RateLimit) (
	awesomeRepos map[string][]*AwesomeRepo, err error) {
	logger := getLogger()
	defer logger.Sync()

	readme, err := client.GetHTMLReadme(awesomeID)
	if err != nil {
		return nil, err
	}
	readmeParser := NewParser(readme, client, reporter, ratelimit)
	awesomeRepos, err = readmeParser.Gather()
	if err != nil {
		return nil, err
	}
	return awesomeRepos, nil
}
