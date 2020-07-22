package main

import (
	"github.com/rydesun/awesome-github/lib/cohttp"
	"github.com/rydesun/awesome-github/lib/errcode"
)

func strerr(err error) string {
	code, scope, _ := errcode.Check(err)
	switch scope {
	case cohttp.ErrScope:
		switch code {
		case cohttp.ErrCodeNetwork:
			return "Network error occurs. Check your network connection."
		}
	}
	return err.Error()
}
