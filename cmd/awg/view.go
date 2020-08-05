package main

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/rydesun/awesome-github/awg"
	"github.com/rydesun/awesome-github/exch/github"
)

type Output struct {
	Time        time.Time                     `json:"time"`
	AwesomeList github.RepoID                 `json:"awesome_list"`
	Data        map[string][]*awg.AwesomeRepo `json:"data"`
	Invalid     []github.RepoID               `json:"invalid"`
}

func LoadOutputFile(path string) (Output, error) {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return Output{}, err
	}
	data := Output{}
	err = json.Unmarshal(raw, &data)
	return data, err
}

func (o *Output) IsValid() bool {
	return len(o.Data) > 0 &&
		!o.Time.IsZero() &&
		len(o.AwesomeList.Name) > 0 &&
		len(o.AwesomeList.Owner) > 0
}
