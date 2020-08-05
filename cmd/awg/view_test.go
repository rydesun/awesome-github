package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOutput(t *testing.T) {
	require := require.New(t)
	datafile := "../../test/testdata/output/data.json"
	data, err := LoadOutputFile(datafile)
	require.Nil(err)
	require.True(data.IsValid())
}
