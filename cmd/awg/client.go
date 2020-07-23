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

	// Check access token and user name.
	fmt.Fprintf(os.Stdout, "[1/3] Checking access token...\n")
	user, err := client.GetUser()
	if err != nil {
		errMsg := "failed to get user info"
		logger.Error(errMsg, zap.Error(err))
		fmt.Println(strerr(err))
		return err
	}
	fmt.Fprintf(os.Stdout, "Use user(%s) access token.\n", user.Name)

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
	fmt.Fprintf(os.Stdout, "[2/3] Parse awesome page...\n")
	awesomeRepos, err := awg.Workflow(client, reporter, config.ID)
	if err != nil {
		errMsg := "failed to fetch awesome repositories"
		logger.Error(errMsg, zap.Error(err))
		fmt.Println(strerr(err))
		return err
	}
	<-finishBar

	// output data
	data, err := json.Marshal(awesomeRepos)
	if err != nil {
		logger.DPanic(err.Error())
		fmt.Println(err.Error())
		return err
	}
	if len(config.Output.Path) != 0 {
		err := ioutil.WriteFile(config.Output.Path, data, 0644)
		if err != nil {
			logger.DPanic(err.Error())
			fmt.Println(err.Error())
		}
	} else {
		fmt.Println(string(data))
	}
	fmt.Printf("Invalid repos: %v\n", reporter.GetInvalidRepo())
	return nil
}
