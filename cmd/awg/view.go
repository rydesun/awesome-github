package main

import (
	"time"

	"github.com/rydesun/awesome-github/awg"
	"github.com/rydesun/awesome-github/exch/github"
)

type Output struct {
	Time    time.Time                     `json:"time"`
	Data    map[string][]*awg.AwesomeRepo `json:"data"`
	Invalid []github.RepoID               `json:"invalid"`
}
