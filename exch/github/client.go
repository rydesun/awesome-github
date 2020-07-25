package github

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"

	"go.uber.org/zap"

	"github.com/rydesun/awesome-github/lib/cohttp"
	"github.com/rydesun/awesome-github/lib/errcode"
)

type HTTPClient interface {
	Json(req *http.Request, resp interface{}) (err error)
	Text(req *http.Request) (text string, err error)
}

type ClientOption struct {
	HTMLHost    string // GitHub HTML host
	HTMLPathPre string // GitHub HTML path prefix
	APIHost     string // GitHub API host
	ApiPathPre  string // GitHub API path prefix
	AccessToken string // GitHub personal access token
}

type Client struct {
	option     ClientOption
	htmlClient HTTPClient
	apiClient  HTTPClient
	htmlHost   url.URL
	apiHost    url.URL
	bearer     string
}

func NewClient(htmlClient HTTPClient, apiClient HTTPClient, option ClientOption) (*Client, error) {
	const funcIntent = "create new github client"
	const funcErrMsg = "failed to " + funcIntent
	logger := getLogger()
	defer logger.Sync()
	logger.Debug(funcIntent, zap.Any("option", option))

	htmlHost, err := url.Parse(option.HTMLHost)
	if err != nil {
		const errMsg = "failed to parse html host"
		logger.Error(errMsg, zap.Error(err),
			zap.String("url", option.HTMLHost))
		err = errcode.New(errMsg, ErrCodeClientOption,
			ErrScope, []string{"htmlHost"})
		return nil, err
	}
	apiHost, err := url.Parse(option.APIHost)
	if err != nil {
		const errMsg = "failed to parse api host"
		logger.Error(errMsg, zap.Error(err),
			zap.String("url", option.APIHost))
		err = errcode.New(errMsg, ErrCodeClientOption,
			ErrScope, []string{"apiHost"})
		return nil, err
	}
	return &Client{
		option:     option,
		htmlClient: htmlClient,
		apiClient:  apiClient,
		htmlHost:   *htmlHost,
		apiHost:    *apiHost,
		bearer:     "Bearer " + option.AccessToken,
	}, nil
}

// Get current user information.
func (c *Client) GetUser() (*User, error) {
	const funcIntent = "get current user info"
	const funcErrMsg = "failed to " + funcIntent
	logger := getLogger()
	defer logger.Sync()

	logger.Debug(funcIntent)

	url := c.apiHost
	url.Path = c.option.ApiPathPre
	req, err := http.NewRequest(http.MethodPost, url.String(),
		bytes.NewBufferString(QueryUser))
	if err != nil {
		errMsg := "failed to create request"
		logger.Error(funcErrMsg, zap.Error(err))
		err = errcode.New(errMsg, ErrCodeMakeRequest,
			ErrScope, nil)
		return nil, err
	}
	req.Header.Set("Authorization", c.bearer)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	user := &User{}
	err = c.apiClient.Json(req, user)
	if err != nil {
		logger.Error(funcErrMsg, zap.Error(err))
		err, ok := err.(errcode.Error)
		if ok && err.Scope == cohttp.ErrScope && err.Code == 401 {
			errMsg := "Invalid access token"
			return nil, errcode.New(errMsg,
				ErrCodeAccessToken, ErrScope, nil)
		}
		return nil, errcode.Wrap(err, funcErrMsg)
	}
	if len(user.Errors) > 0 {
		errMsg := "remote server return errors"
		logger.Error(errMsg, zap.Any("serverErrors", user.Errors))
		err = errcode.New(errMsg, errcode.CodeUnknown, ErrScope, nil)
		return nil, err
	}
	return user, nil
}

// Get repository information.
func (c Client) GetRepo(id RepoID) (*Repo, error) {
	const funcIntent = "get repo info"
	const funcErrMsg = "failed to " + funcIntent
	logger := getLogger()
	defer logger.Sync()

	idStr := id.String()
	logger.Debug(funcIntent, zap.String("repo", idStr))

	url := c.apiHost
	url.Path = c.option.ApiPathPre

	query := fmt.Sprintf(QueryRepo, id.Owner, id.Name)
	req, err := http.NewRequest(http.MethodPost, url.String(),
		bytes.NewBufferString(query))
	if err != nil {
		errMsg := "failed to create request"
		logger.Error(funcErrMsg, zap.Error(err),
			zap.String("repo", idStr))
		err = errcode.New(errMsg, ErrCodeMakeRequest,
			ErrScope, nil)
		return nil, err
	}
	req.Header.Set("Authorization", c.bearer)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	repo := &Repo{}
	err = c.apiClient.Json(req, repo)
	if err != nil {
		logger.Error(funcErrMsg, zap.Error(err),
			zap.String("repo", idStr))
		err = errcode.Wrap(err, funcErrMsg)
		return nil, err
	}
	if len(repo.Errors) > 0 {
		errMsg := "remote server return errors"
		logger.Error(errMsg, zap.String("repo", idStr),
			zap.Any("serverErrors", repo.Errors))
		err = errcode.New(errMsg, errcode.CodeUnknown, ErrScope, nil)
		return nil, err
	}
	return repo, nil
}

// Get README.md wrapped in html page.
func (c Client) GetHTMLReadme(id RepoID) (string, error) {
	const funcIntent = "get repo readme wrapped in html page"
	const funcErrMsg = "failed to " + funcIntent
	logger := getLogger()
	defer logger.Sync()

	idStr := id.String()
	logger.Debug(funcIntent, zap.String("repo", idStr))

	readme, err := c.GetHTML(id, "README.md")
	if err != nil {
		logger.Error(funcErrMsg, zap.Error(err),
			zap.String("repo", idStr))
		err = errcode.Wrap(err, funcErrMsg)
		return "", err
	}
	return readme, nil
}

// Get file wrapped in html page.
func (c *Client) GetHTML(id RepoID, path string) (string, error) {
	const funcIntent = "get file wrapped in html page"
	const funcErrMsg = "failed to " + funcIntent
	logger := getLogger()
	defer logger.Sync()

	idStr := id.String()
	logger.Debug(funcIntent,
		zap.String("repo", idStr),
		zap.String("path", path))

	url := c.htmlHost
	url.Path = filepath.Join(idStr, c.option.HTMLPathPre, path)
	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		errMsg := "failed to create request"
		logger.Error(errMsg, zap.Error(err),
			zap.String("repo", idStr),
			zap.String("path", path))
		err = errcode.New(errMsg, ErrCodeMakeRequest,
			ErrScope, nil)
		return "", err
	}
	content, err := c.htmlClient.Text(req)
	if err != nil {
		logger.Error(funcErrMsg, zap.Error(err),
			zap.String("repo", idStr),
			zap.String("path", path))
		err = errcode.New(funcErrMsg, ErrCodeNetwork,
			ErrScope, []string{"github"})
		return "", err
	}
	return content, nil
}
