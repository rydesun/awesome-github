package awg

import (
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/antchfx/htmlquery"
	"go.uber.org/zap"
	"golang.org/x/net/html"

	"github.com/rydesun/awesome-github/exch/github"
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
}

func NewParser(readme string, client *Client, reporter *Reporter) (
	*Parser, error) {
	const funcIntent = "parse awesome html readme page"
	const funcErrMsg = "failed to " + funcIntent
	logger := getLogger()
	defer logger.Sync()

	logger.Debug(funcIntent)

	node, err := htmlquery.Parse(strings.NewReader(readme))
	if err != nil {
		errMsg := "failed to recognize readme html page content"
		logger.Error(errMsg, zap.Error(err))
		err = errcode.New(errMsg, ErrCodeContent,
			ErrScope, []string{"readme"})
		return nil, err
	}
	return &Parser{
		client:       client,
		node:         node,
		xpathSection: xpathSection,
		xpathItem:    xpathItem,
		urlMust:      urlMust,
		reLink:       reLink,
		reporter:     reporter,
	}, nil
}

// Gather repositories from awesome README.md.
func (p *Parser) Gather() (map[string][]*AwesomeRepo, error) {
	const funcIntent = "gather repositories"
	const funcErrMsg = "failed to " + funcIntent
	logger := getLogger()
	defer logger.Sync()
	logger.Debug(funcIntent)

	sectionNodes, err := p.getSections(p.xpathSection)
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
			err = nil
			continue
		}
		sectionItemsMap[sectionName] = itemNodes
	}

	// Section -> Index
	// Item -> Repo
	idxReposMap := make(map[string][]*Repo, len(sectionItemsMap))
	for sectionName, itemNodes := range sectionItemsMap {
		repos := make([]*Repo, 0)
		for _, itemNode := range itemNodes {
			repo, err := p.parseItem(itemNode)
			if err != nil {
				logger.Warn("skip invalid item", zap.Error(err))
				err = nil
				continue
			}
			repos = append(repos, repo)
		}
		idxReposMap[sectionName] = repos
	}

	var wg sync.WaitGroup
	var jobNum int
	idxAwReposMap := make(map[string][]*AwesomeRepo, len(idxReposMap))
	for idx, repos := range idxReposMap {
		for cnt, repo := range repos {
			awesomeDesc := p.getDesc(sectionItemsMap[idx][cnt])
			awesomeRepo := &AwesomeRepo{
				Repo:        *repo,
				AwesomeDesc: awesomeDesc,
			}
			idxAwReposMap[idx] = append(idxAwReposMap[idx], awesomeRepo)
			wg.Add(1)
			jobNum++
			go func(idx string, cnt int) {
				defer wg.Done()
				err := p.client.Fill(awesomeRepo)
				if p.reporter != nil {
					p.reporter.Done()
				}
				if err != nil {
					errMsg := "failed to fill repository info"
					logger.Error(errMsg, zap.Error(err))
					if p.reporter != nil {
						p.reporter.InvalidRepo(awesomeRepo.ID)
					}
					idxAwReposMap[idx][cnt] = nil
				}
			}(idx, cnt)
		}
	}
	if p.reporter != nil {
		p.reporter.TotalRepoNum(jobNum)
	}
	wg.Wait()
	return p.clean(idxAwReposMap), nil
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

// Get awesome section nodes from awesome README.md
func (p *Parser) getSections(xpath string) ([]*html.Node, error) {
	return htmlquery.QueryAll(p.node, xpath)
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

	linkNode, err := htmlquery.Query(itemNode, "//a")
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
