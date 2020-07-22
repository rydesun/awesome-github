package errcode

type ErrScope string

const (
	ScopeUnknown  ErrScope = "unknown"
	ScopeInternal ErrScope = "internal"
)

type ErrCode int

const (
	CodeUnknown  ErrCode = 0
	CodeInternal ErrCode = 1
)
