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
	Watch       int           `json:"watch"`
	Star        int           `json:"star"`
	Fork        int           `json:"fork"`
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
	r.Watch = repo.Data.Repository.Watchers.TotalCount
	r.Star = repo.Data.Repository.Stargazers.TotalCount
	r.Fork = repo.Data.Repository.Forks.TotalCount
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
