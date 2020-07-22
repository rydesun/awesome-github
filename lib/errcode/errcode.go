package errcode

import (
	"encoding/json"

	"go.uber.org/zap"
)

type Error struct {
	Msg     string
	Code    ErrCode
	Scope   ErrScope
	Obejcts []string
	err     string
}

func (e Error) Error() string {
	raw, _ := json.Marshal(e)
	return string(raw)
}

func New(msg string, code ErrCode, scope ErrScope, obejcts []string) error {
	if len(scope) == 0 {
		scope = ScopeUnknown
	}
	if obejcts == nil {
		obejcts = []string{}
	}
	return Error{
		Msg:     msg,
		Code:    code,
		Scope:   scope,
		Obejcts: obejcts,
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

func Check(err error) (errCode ErrCode, errScope ErrScope, objects []string) {
	e, ok := err.(Error)
	if ok {
		return e.Code, e.Scope, e.Obejcts
	} else {
		return CodeUnknown, ScopeUnknown, []string{}
	}
}
