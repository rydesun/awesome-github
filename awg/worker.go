package awg

import (
	"github.com/rydesun/awesome-github/exch/github"
)

func Workflow(client *Client, reporter *Reporter, awesomeID github.RepoID) (
	awesomeRepos map[string][]*AwesomeRepo, err error) {
	logger := getLogger()
	defer logger.Sync()

	readme, err := client.GetHTMLReadme(awesomeID)
	if err != nil {
		return nil, err
	}
	readmeParser, err := NewParser(readme, client, reporter)
	if err != nil {
		return nil, err
	}
	awesomeRepos, err = readmeParser.Gather()
	if err != nil {
		return nil, err
	}
	return awesomeRepos, nil
}
