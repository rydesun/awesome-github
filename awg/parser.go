package awg

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/antchfx/htmlquery"
	"go.uber.org/zap"
	"golang.org/x/net/html"

	"github.com/rydesun/awesome-github/exch/github"
	"github.com/rydesun/awesome-github/lib/cohttp"
	"github.com/rydesun/awesome-github/lib/errcode"
)

type Parser struct {
	client       *Client
	node         *html.Node
	xpathSection string
	xpathItem    string
	urlMust      string
	reLink       *regexp.Regexp
	reporter     *Reporter
	ratelimit    RateLimit
}

func NewParser(readme string, client *Client, reporter *Reporter,
	rateLimit RateLimit) *Parser {
	// Don't worry about this error.
	node, _ := htmlquery.Parse(strings.NewReader(readme))
	return &Parser{
		client:       client,
		node:         node,
		xpathSection: xpathSection,
		xpathItem:    xpathItem,
		urlMust:      urlMust,
		reLink:       reLink,
		reporter:     reporter,
		ratelimit:    rateLimit,
	}
}

// Gather repositories from awesome README.md.
func (p *Parser) Gather() (map[string][]*AwesomeRepo, error) {
	const funcIntent = "gather repositories"
	const funcErrMsg = "failed to " + funcIntent
	logger := getLogger()
	defer logger.Sync()
	logger.Debug(funcIntent)

	sectionItemsMap, err := p.Parse()
	if err != nil {
		return nil, err
	}
	if len(sectionItemsMap) == 0 {
		errMsg := "failed to find any valid sections"
		logger.Error(errMsg)
		return nil, errcode.New(errMsg, ErrCodeContent, ErrScope, nil)
	}
	total, idxReposMap := p.convert(sectionItemsMap)
	// TODO: may be different with graphQL node number
	if total > p.ratelimit.Remaining {
		errMsg := "Exceed GitHub API ratelimit"
		logger.Warn(errMsg, zap.Error(err))
		return nil, errcode.New(errMsg, ErrCodeRatelimit, ErrScope,
			[]string{strconv.Itoa(total)})
	}
	if p.reporter != nil {
		p.reporter.TotalRepoNum(total)
	}

	err = p.FetchRepos(idxReposMap)
	if err != nil {
		return nil, err
	}
	idxReposMap = p.clean(idxReposMap)
	return idxReposMap, nil
}

// Get awesome section nodes from awesome README.md
func (p *Parser) Parse() (map[string][]*html.Node, error) {
	logger := getLogger()
	defer logger.Sync()

	sectionNodes, err := htmlquery.QueryAll(p.node, p.xpathSection)
	if len(sectionNodes) == 0 {
		errMsg := "awesome html page does not contain any sections"
		logger.Error(errMsg, zap.Error(err))
		return nil, errcode.New(errMsg, ErrCodeContent,
			ErrScope, []string{"section"})
	}
	logger.Info("get some section nodes", zap.Int("len", len(sectionNodes)))

	sectionItemsMap := make(map[string][]*html.Node, 0)
	for i, sectionNode := range sectionNodes {
		sectionName := htmlquery.InnerText(sectionNode)
		itemNodes, err := p.getItemsFromSection(sectionNode, i, p.xpathItem)
		if err != nil || len(itemNodes) == 0 {
			errMsg := "wired section has no items"
			logger.Warn(errMsg, zap.Error(err),
				zap.String("section", sectionName))
			continue
		}
		sectionItemsMap[sectionName] = itemNodes
	}
	return sectionItemsMap, nil
}

// Section -> Index
// Item -> AwesomeRepo
func (p *Parser) convert(sectionItemsMap map[string][]*html.Node) (
	total int, idxReposMap map[string][]*AwesomeRepo) {
	logger := getLogger()
	defer logger.Sync()
	idxReposMap = make(map[string][]*AwesomeRepo, len(sectionItemsMap))

	for sectionName, itemNodes := range sectionItemsMap {
		repos := make([]*AwesomeRepo, 0)
		for _, itemNode := range itemNodes {
			repo, err := p.parseItem(itemNode)
			if err != nil {
				logger.Warn("skip invalid item", zap.Error(err))
				continue
			}
			repos = append(repos, &AwesomeRepo{
				Repo:        *repo,
				AwesomeDesc: p.getDesc(itemNode),
			})
			total++
		}
		idxReposMap[sectionName] = repos
	}
	return total, idxReposMap
}

// Fetch repositories from remote.
func (p *Parser) FetchRepos(idxReposMap map[string][]*AwesomeRepo) error {
	logger := getLogger()
	defer logger.Sync()

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	unacceptedError := make(chan error)
	for idx, repos := range idxReposMap {
		for cnt, repo := range repos {
			wg.Add(1)
			go func(repo *AwesomeRepo, idx string, cnt int) {
				defer wg.Done()
				err := p.client.Fill(ctx, repo)

				if p.reporter != nil {
					p.reporter.Done()
				}
				if err != nil {
					if cohttp.IsNetowrkError(err) {
						errMsg := "Network error occurs"
						err = errcode.Wrap(err, errMsg)
						unacceptedError <- err
						return
					} else if github.IsAbuseError(err) {
						errMsg := "The frequency of requests is too high. Check max_concurrent"
						err = errcode.Wrap(err, errMsg)
						unacceptedError <- err
						return
					}
					// accepted error
					errMsg := "failed to fill repository info"
					logger.Error(errMsg, zap.Error(err))
					if p.reporter != nil {
						p.reporter.InvalidRepo(repo.ID)
					}
					idxReposMap[idx][cnt] = nil
				}
			}(repo, idx, cnt)
		}
	}
	jobsCompleted := make(chan struct{})
	go func() {
		wg.Wait()
		close(jobsCompleted)
	}()
	select {
	case err := <-unacceptedError:
		return err
	case <-jobsCompleted:
		return nil
	}
}

// Remove invalid nil from map.
func (p *Parser) clean(raw map[string][]*AwesomeRepo) map[string][]*AwesomeRepo {
	result := make(map[string][]*AwesomeRepo, len(raw))

	for idx, repos := range raw {
		for _, repo := range repos {
			if repo != nil {
				result[idx] = append(result[idx], repo)
			}
		}
	}
	return result
}

// Get awesome item nodes from awesome section.
func (p *Parser) getItemsFromSection(section *html.Node, idx int, xpath string) (
	[]*html.Node, error) {
	xpath = fmt.Sprintf(xpath, idx+1)
	return htmlquery.QueryAll(section, xpath)
}

// Get awesome link node from a item node.
func (p *Parser) getLinks(itemNode *html.Node) (*html.Node, error) {
	const funcIntent = "get awesome links from a section node"
	const funcErrMsg = "failed to " + funcIntent
	logger := getLogger()
	defer logger.Sync()
	logger.Debug(funcIntent)

	return htmlquery.Query(itemNode, "//a")
}

// Get awesome description from a item node.
func (p *Parser) getDesc(itemNode *html.Node) string {
	return htmlquery.InnerText(itemNode)
}

// Get repository from a item node.
// One item one repo.
func (p *Parser) parseItem(itemNode *html.Node) (*Repo, error) {
	const funcIntent = "get repo from a item node"
	const funcErrMsg = "failed to " + funcIntent
	logger := getLogger()
	defer logger.Sync()
	logger.Debug(funcIntent)

	linkNode, err := p.getLinks(itemNode)
	if err != nil {
		const blockErrMsg = "failed to get first link node from item node"
		logger.Error(blockErrMsg, zap.Error(err))
		return nil, errcode.Wrap(err, blockErrMsg)
	}

	name, link :=
		htmlquery.InnerText(linkNode),
		htmlquery.SelectAttr(linkNode, "href")
	if name == "" || link == "" {
		logger.Error(funcErrMsg)
		return nil, errcode.New(funcErrMsg, ErrCodeContent,
			ErrScope, nil)
	}

	linkSplit := p.reLink.FindStringSubmatch(link)
	if len(linkSplit) != 3 || strings.Contains(linkSplit[2], "/") {
		var errMsg string
		if strings.Contains(link, p.urlMust) {
			errMsg = "strange github repository url"
			logger.Warn(errMsg,
				zap.String("link", link),
				zap.String("name", name))
		} else {
			errMsg = "discard unrecognized url"
			logger.Info(errMsg,
				zap.String("link", link),
				zap.String("name", name))
		}
		return nil, errcode.New(errMsg, ErrCodeContent, ErrScope, nil)
	}

	id := github.RepoID{
		Owner: linkSplit[1],
		Name:  linkSplit[2],
	}
	return &Repo{
		ID:          id,
		AwesomeName: name,
		Owner:       id.Owner,
		Link:        link,
	}, nil
}
