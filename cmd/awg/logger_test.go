package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/rydesun/awesome-github/exch/config"
	"github.com/stretchr/testify/require"
)

func TestSetLoggers(t *testing.T) {
	require := require.New(t)
	tempFile, err := ioutil.TempFile("", "awg_fake_log_*")
	require.Nil(err)
	defer func() {
		tempFile.Close()
		os.Remove(tempFile.Name())
	}()
	mainLogger, err := setLoggers(config.Loggers{
		Main: config.Logger{
			Path: tempFile.Name(),
		},
		Http: config.Logger{
			Path: tempFile.Name(),
		},
	})
	require.Nil(err)
	mainLogger.Info("test log")
	err = mainLogger.Sync()
	require.Nil(err)
}
