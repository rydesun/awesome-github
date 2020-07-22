package errcode

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheck(t *testing.T) {
	assert := assert.New(t)

	expectedCode := CodeInternal
	expectedScope := ErrScope("test")
	expectedObjects := []string{"a", "b"}

	err := New("test", expectedCode, ErrScope(expectedScope), expectedObjects)
	code, scope, objects := Check(err)
	assert.Equal(expectedCode, code)
	assert.Equal(expectedScope, scope)
	assert.Equal(expectedObjects, objects)

	err = Wrap(err, "wrap")
	code, scope, objects = Check(err)
	assert.Equal(expectedCode, code)
	assert.Equal(expectedScope, scope)
	assert.Equal(expectedObjects, objects)

	err = New("test", 0, "", nil)
	code, scope, objects = Check(err)
	assert.Equal(ErrCode(0), code)
	assert.Equal(ErrScope("unknown"), scope)
	assert.Equal([]string{}, objects)

	abnormalErr := errors.New("simple error")

	err = Wrap(abnormalErr, "wrap")
	code, scope, objects = Check(err)
	assert.Equal(CodeUnknown, code)
	assert.Equal(ScopeUnknown, scope)
	assert.Equal([]string{}, objects)

	code, scope, objects = Check(abnormalErr)
	assert.Equal(CodeUnknown, code)
	assert.Equal(ScopeUnknown, scope)
	assert.Equal([]string{}, objects)
}
