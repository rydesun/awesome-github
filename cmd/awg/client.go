package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
	"go.uber.org/zap"

	"github.com/rydesun/awesome-github/awg"
	"github.com/rydesun/awesome-github/exch/config"
	"github.com/rydesun/awesome-github/exch/github"
	"github.com/rydesun/awesome-github/lib/cohttp"
)

func newReporter() *awg.Reporter {
	return &awg.Reporter{}
}

func newClient(config config.Config, reporter *awg.Reporter) (
	*awg.Client, error) {
	logger := getLogger()
	defer logger.Sync()

	htmlClient := cohttp.NewClient(http.Client{},
		config.MaxConcurrent, config.LogRespHead, nil)
	apiClient := cohttp.NewClient(http.Client{},
		config.MaxConcurrent, config.LogRespHead, reporter)
	options := github.NewDefaultClientOption()
	options.AccessToken = config.AccessToken
	gc, err := github.NewClient(htmlClient, apiClient, options)
	if err != nil {
		errMsg := "failed to create github client"
		logger.Error(errMsg, zap.Error(err))
		return nil, err
	}
	client, err := awg.NewClient(gc)
	if err != nil {
		errMsg := "failed to create awg client"
		logger.Error(errMsg, zap.Error(err))
		return nil, err
	}
	return client, nil
}

func work(client *awg.Client, config config.Config, reporter *awg.Reporter) error {
	logger := getLogger()
	defer logger.Sync()
	logger.Info("awesome analysis instance begins")

	writer := os.Stdout

	// Check access token and user name.
	fmt.Fprintf(writer, "[1/3] Checking access token...\n")
	user, err := client.GetUser()
	if err != nil {
		errMsg := "failed to get user info"
		logger.Error(errMsg, zap.Error(err))
		fmt.Fprintln(writer, strerr(err))
		return err
	}
	fmt.Fprintf(writer, "Use user(%s) access token.\n", user.Name)
	fmt.Fprintf(writer, "RateLimit: total %d, remaining %d, reset at %s\n",
		user.RateLimitTotal, user.RateLimitRemaining, user.RateLimitResetAt)

	finishBar := make(chan interface{})
	go func() {
		var jobNum int
		for {
			jobNum = reporter.GetTotalRepoNum()
			if jobNum > 0 {
				break
			}
			time.Sleep(time.Second)
		}
		bar := progressbar.NewOptions(jobNum,
			progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
			progressbar.OptionEnableColorCodes(true),
			progressbar.OptionSetWidth(15),
			progressbar.OptionSetDescription("[3/3] Fetch github repo info..."),
			progressbar.OptionShowIts(),
			progressbar.OptionShowCount(),
		)
		for {
			finishedNum := reporter.GetFinishedRepoNum()
			if finishedNum == jobNum {
				bar.Finish()
				break
			}
			bar.Set(finishedNum)
			time.Sleep(time.Second)
		}
		finishBar <- nil
	}()

	// actual work
	fmt.Fprintf(writer, "[2/3] Parse awesome page...\n")
	awesomeRepos, err := awg.Workflow(client, reporter, config.ID)
	if err != nil {
		errMsg := "failed to fetch awesome repositories"
		logger.Error(errMsg, zap.Error(err))
		fmt.Fprintln(writer, strerr(err))
		return err
	}
	<-finishBar

	// output data
	data, err := json.Marshal(awesomeRepos)
	if err != nil {
		logger.DPanic(err.Error())
		fmt.Fprintln(writer, strerr(err))
		return err
	}
	if len(config.Output.Path) != 0 {
		err := ioutil.WriteFile(config.Output.Path, data, 0644)
		if err != nil {
			logger.DPanic(err.Error())
			fmt.Fprintln(writer, strerr(err))
		}
	} else {
		fmt.Fprintln(writer, string(data))
	}

	// Warn some invalid repos.
	invalidRepos := reporter.GetInvalidRepo()
	if len(invalidRepos) > 0 {
		fmt.Fprintf(writer, "\nCatch some invalid repos: %v\n", invalidRepos)
	}

	// last message
	if len(config.Output.Path) > 0 {
		fmt.Fprintf(writer, "Done. Output file: %s\n", config.Output.Path)
	} else {
		fmt.Fprintln(writer, "Done.")
	}
	return nil
}
