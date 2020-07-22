package cohttp

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTruncate(t *testing.T) {
	require := require.New(t)
	testCases := []struct {
		maxLength int
		raw       []byte
		result    []byte
	}{
		{
			maxLength: 0,
			raw:       []byte("abcdefghi"),
			result:    []byte{},
		},
		{
			maxLength: 1,
			raw:       []byte("abcdefghi"),
			result:    []byte("a..."),
		},
		{
			maxLength: 8,
			raw:       []byte("abcdefghi"),
			result:    []byte("abcdefgh..."),
		},
		{
			maxLength: 9,
			raw:       []byte("abcdefghi"),
			result:    []byte("abcdefghi"),
		},
		{
			maxLength: 10,
			raw:       []byte("abcdefghi"),
			result:    []byte("abcdefghi"),
		},
	}
	for _, tc := range testCases {
		t.Run("NONAME", func(t *testing.T) {
			actual := truncate(tc.raw, tc.maxLength)
			require.Equal(tc.result, actual)
		})
	}
}
