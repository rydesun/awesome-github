package awg

import (
	"time"

	"github.com/rydesun/awesome-github/exch/github"
	"github.com/rydesun/awesome-github/lib/errcode"
)

type Repo struct {
	ID          github.RepoID `json:"id"`
	Owner       string        `json:"owner"`
	AwesomeName string        `json:"awesome_name"`
	Link        string        `json:"link"`
	Star        int           `json:"star"`
	LastCommit  time.Time     `json:"last_commit"`
	Description string        `json:"description"`
}

type AwesomeRepo struct {
	Repo
	AwesomeDesc string `json:"awesome_description"`
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

type RateLimit struct {
	Total     int
	Remaining int
	ResetAt   time.Time
}

type User struct {
	Name string
	RateLimit
}
