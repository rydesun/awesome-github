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
	errc, ok := err.(errcode.Error)
	if !ok {
		return err.Error()
	}
	switch errc.Scope {
	case awg.ErrScope:
		switch errc.Code {
		case awg.ErrCodeRatelimit:
			return "Exceed GitHub API ratelimit."
		}
	case github.ErrScope:
		switch errc.Code {
		case github.ErrCodeAccessToken:
			return "Invalid github personal access token."
		}
	case config.ErrScope:
		switch errc.Code {
		case config.ErrCodeParameter:
			return fmt.Sprintf("Invalid config: %v", errc.Msg)
		}
	case cohttp.ErrScope:
		switch errc.Code {
		case cohttp.ErrCodeNetwork:
			msg := "Network error occurs. Check your network connection."
			if len(errc.Objects) == 0 {
				return msg
			}
			msg = msg + "\n" + errc.Objects[0]
			return msg
		}
	}
	return err.Error()
}
