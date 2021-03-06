package awg

import (
	"context"

	"go.uber.org/zap"

	"github.com/rydesun/awesome-github/exch/github"
	"github.com/rydesun/awesome-github/lib/errcode"
)

type Client struct {
	gc *github.Client
}

// New awg client.
func NewClient(client *github.Client) (*Client, error) {
	return &Client{
		gc: client,
	}, nil
}

func (c *Client) GetUser() (*User, error) {
	user, err := c.gc.GetUser()
	if err != nil {
		return nil, err
	}
	return &User{
		Name: user.Data.Viewer.Login,
		RateLimit: RateLimit{
			Total:     user.Data.Ratelimit.Limit,
			Remaining: user.Data.Ratelimit.Remaining,
			ResetAt:   user.Data.Ratelimit.ResetAt,
		},
	}, nil
}

// Get Readme html page.
func (c *Client) GetHTMLReadme(id github.RepoID) (string, error) {
	const funcIntent = "get readme html page"
	const funcErrMsg = "failed to " + funcIntent
	return c.gc.GetHTMLReadme(id)
}

// Fill struct repo with more info.
func (c *Client) Fill(ctx context.Context, repo *AwesomeRepo) error {
	const funcIntent = "fill struct repo with more info"
	const funcErrMsg = "failed to " + funcIntent
	logger := getLogger()
	defer logger.Sync()

	id := repo.ID
	idStr := repo.ID.String()

	logger.Debug(funcIntent, zap.String("repo", idStr))

	rawRepo, err := c.gc.GetRepo(ctx, id)
	if err != nil {
		logger.Error(funcErrMsg, zap.Error(err),
			zap.String("repo", idStr))
		return errcode.Wrap(err, funcErrMsg)
	}
	repo.Aggregate(rawRepo)
	return nil
}
