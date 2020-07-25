package main

import (
	"fmt"

	"github.com/rydesun/awesome-github/awg"
	"github.com/rydesun/awesome-github/exch/config"
	"github.com/rydesun/awesome-github/exch/github"
	"github.com/rydesun/awesome-github/lib/cohttp"
	"github.com/rydesun/awesome-github/lib/errcode"
)

func strerr(err error) string {
	code, scope, objects := errcode.Check(err)
	switch scope {
	case awg.ErrScope:
		switch code {
		case awg.ErrCodeRatelimit:
			return "Exceed GitHub API ratelimit."
		}
	case github.ErrScope:
		switch code {
		case github.ErrCodeAccessToken:
			return "Invalid github personal access token."
		}
	case config.ErrScope:
		switch code {
		case config.ErrCodeParameter:
			return fmt.Sprintf("Invalid config: %v", err)
		}
	case cohttp.ErrScope:
		switch code {
		case cohttp.ErrCodeNetwork:
			msg := "Network error occurs. Check your network connection."
			if len(objects) == 0 {
				return msg
			}
			msg = msg + "\n" + objects[0]
			return msg
		}
	}
	return err.Error()
}
