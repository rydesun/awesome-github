package awg

import "regexp"

const (
	ErrScope = "awg"

	ErrCodeContent = 10
)

const xpathSection = `//div[@id='readme']
	//article[contains(@class,'markdown-body')]
	//h2`
const xpathItem = "//following-sibling::ul[count(preceding-sibling::h2)=%v]/li"
const urlMust = "github.com"

var reLink = regexp.MustCompile("(?U)^https?://github.com/(.+)/(.+)/?$")
