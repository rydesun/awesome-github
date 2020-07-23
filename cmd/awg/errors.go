package main

import (
	"github.com/rydesun/awesome-github/awg"
	"github.com/rydesun/awesome-github/lib/cohttp"
	"github.com/rydesun/awesome-github/lib/errcode"
)

func strerr(err error) string {
	code, scope, _ := errcode.Check(err)
	switch scope {
	case awg.ErrScope:
		switch code {
		case awg.ErrCodeRatelimit:
			return "Exceed GitHub API ratelimit."
		}
	case cohttp.ErrScope:
		switch code {
		case cohttp.ErrCodeNetwork:
			return "Network error occurs. Check your network connection."
		}
	}
	return err.Error()
}
