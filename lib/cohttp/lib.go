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

func IsNetowrkError(err error) bool {
	errc, ok := err.(errcode.Error)
	return ok && errc.Code == ErrCodeNetwork
}
