package errcode

import (
	"encoding/json"

	"go.uber.org/zap"
)

type Error struct {
	Msg     string
	Code    ErrCode
	Scope   ErrScope
	Objects []string
	err     string
}

func (e Error) Error() string {
	raw, _ := json.Marshal(e)
	return string(raw)
}

func New(msg string, code ErrCode, scope ErrScope, objects []string) error {
	if len(scope) == 0 {
		scope = ScopeUnknown
	}
	if objects == nil {
		objects = []string{}
	}
	return Error{
		Msg:     msg,
		Code:    code,
		Scope:   scope,
		Objects: objects,
	}
}

func Wrap(err error, msg string) error {
	e, ok := err.(Error)
	if ok {
		e.Msg = msg
		return e
	}
	logger := getLogger()
	defer logger.Sync()
	logger.DPanic("generate a unknown error",
		zap.Error(err))
	return New(msg, CodeUnknown, ScopeUnknown, []string{})
}
