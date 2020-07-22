package github

const (
	ErrScope = "github"

	ErrCodeParameter    = 10
	ErrCodeClientOption = 11
	ErrCodeMakeRequest  = 20
	ErrCodeNetwork      = 30
)

const (
	HTMLHost    = "https://github.com/"
	HTMLPathPre = "blob/master"
	APIHost     = "https://api.github.com/"
	APIPathPre  = "graphql"
)

func NewDefaultClientOption() ClientOption {
	return ClientOption{
		HTMLHost:    HTMLHost,
		HTMLPathPre: HTMLPathPre,
		APIHost:     APIHost,
		ApiPathPre:  APIPathPre,
	}
}

const QueryRepo = `{ "query": "query { repository(owner: \"%s\", name: \"%s\") { description forks { totalCount } stargazers { totalCount } watchers { totalCount } defaultBranchRef { target { ... on Commit { history(first: 1) { edges { node { committedDate } } } } } } } }" }`

const QueryUser = `{ "query": "query { viewer { login }}"`
