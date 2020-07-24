package cohttp

import "github.com/rydesun/awesome-github/lib/errcode"

func truncate(raw []byte, maxLength int) (result []byte) {
	if maxLength == 0 {
		return []byte{}
	}
	length := len(raw)
	if length > maxLength {
		return append(raw[:maxLength], "..."...)
	} else {
		return raw
	}
}

func IsNetowrkError(err error) (int, bool) {
	code, _, _ := errcode.Check(err)
	if code < 100 {
		return int(code), true
	} else {
		return 0, false
	}
}
