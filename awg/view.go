package awg

import (
	"time"

	"github.com/rydesun/awesome-github/exch/github"
	"github.com/rydesun/awesome-github/lib/errcode"
)

type Repo struct {
	ID          github.RepoID
	Owner       string
	AwesomeName string
	Link        string
	Star        int
	LastCommit  time.Time
	Description string
}

type AwesomeRepo struct {
	Repo
	AwesomeDesc string
}

func (r *Repo) Aggregate(repo *github.Repo) error {
	commitEdges := repo.Data.Repository.DefaultBranchRef.Target.History.Edges
	if len(commitEdges) == 0 {
		errMsg := "malformed repo struct"
		return errcode.New(errMsg, ErrCodeContent,
			ErrScope, []string{"commit"})
	}
	r.LastCommit = commitEdges[0].Node.CommittedDate
	r.Star = repo.Data.Repository.Stargazers.TotalCount
	r.Description = repo.Data.Repository.Description
	return nil
}

type User struct {
	Name               string
	RateLimitTotal     int
	RateLimitRemaining int
	RateLimitResetAt   time.Time
}
