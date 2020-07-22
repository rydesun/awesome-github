package github

import (
	"fmt"
	"time"
)

type RepoID struct {
	Owner string
	Name  string
}

func (r RepoID) String() string {
	return fmt.Sprintf("%s/%s", r.Owner, r.Name)
}

type Repo struct {
	Data struct {
		Repository struct {
			Forks struct {
				TotalCount int
			}
			Stargazers struct {
				TotalCount int
			}
			Watchers struct {
				TotalCount int
			}
			DefaultBranchRef struct {
				Target struct {
					History struct {
						Edges []struct {
							Node struct {
								CommittedDate time.Time
							}
						}
					}
				}
			}
			Description string
		}
	}
	Errors []struct {
		Message string
	}
}

type User struct {
	Data struct {
		Viewer struct {
			Login string
		}
	}
	Errors []struct {
		Message string
	}
}
