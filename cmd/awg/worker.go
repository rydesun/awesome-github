package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
	"go.uber.org/zap"

	"github.com/rydesun/awesome-github/awg"
	"github.com/rydesun/awesome-github/exch/config"
	"github.com/rydesun/awesome-github/exch/github"
	"github.com/rydesun/awesome-github/lib/cohttp"
)

type Worker struct {
	repoID     github.RepoID
	outputPath string
	writer     io.Writer
	reporter   *awg.Reporter
	awgClient  *awg.Client
	logger     *zap.Logger

	// CLI settings
	disableProgressBar bool
}

func NewWorker(writer io.Writer, logger *zap.Logger) *Worker {
	return &Worker{
		writer: writer,
		logger: logger,
	}
}

func (w *Worker) Init(config config.Config) error {
	writer := w.writer

	// Introdution.
	fmt.Fprintln(writer, "=== Awesome GitHub ===")

	// Show config.
	fmt.Fprintf(writer, "config file: %s\n", config.ConfigPath)
	fmt.Fprintf(writer, "main log files: %s\n", config.Log.Main.Path)
	fmt.Fprintf(writer, "http log files: %s\n", config.Log.Http.Path)
	fmt.Fprintf(writer, "output file: %s\n", config.Output.Path)

	// Create awg client.
	w.reporter = w.newReporter()
	awgClient, err := w.newAwgClient(config)
	if err != nil {
		fmt.Fprintln(writer, err)
		return err
	}

	w.repoID = config.ID
	w.outputPath = config.Output.Path
	w.awgClient = awgClient
	w.disableProgressBar = config.Cli.DisableProgressBar
	return nil
}

func (w *Worker) Work() error {
	writer := w.writer
	logger := w.logger
	defer logger.Sync()

	// Check access token.
	fmt.Fprintf(w.writer, "[1/3] Checking github access token...\n")
	user, err := w.awgClient.GetUser()
	if err != nil {
		errMsg := "Failed to get information about access token."
		fmt.Fprintln(w.writer, errMsg, strerr(err))
		logger.Error(errMsg, zap.Error(err))
		return err
	}
	fmt.Fprintf(w.writer, "Use user(%s) access token.\n", user.Name)
	fmt.Fprintf(w.writer, "RateLimit: total %d, remaining %d, reset at %s\n",
		user.RateLimit.Total, user.RateLimit.Remaining, user.RateLimit.ResetAt)

	fmt.Fprintln(writer, "[2/3] Fetch and parse awesome README.md...")

	// Progress bar.
	var pbCompleted <-chan struct{}
	var pbCancel chan<- struct{}
	if !w.disableProgressBar {
		pbCompleted, pbCancel = w.progressBar("[3/3] Fetch repositories from github...")
	} else {
		fmt.Fprintln(writer, "[3/3] Fetch repositories from github...")
	}

	// Actual work.
	awesomeRepos, err := awg.Workflow(w.awgClient, w.reporter, w.repoID, user.RateLimit)
	if err != nil {
		if !w.disableProgressBar {
			logger.Info("Cancel progress bar.")
			close(pbCancel)
			logger.Info("progress bar canceled.")
		}
		errMsg := "\nFailed to fetch some repositories."
		fmt.Fprintln(writer, errMsg, strerr(err))
		logger.Error(errMsg, zap.Error(err))
		return err
	}
	if !w.disableProgressBar {
		logger.Info("Wait for the progress bar to complete.")
		<-pbCompleted
		logger.Info("Progress bar finished.")
	}
	invalidRepos := w.reporter.GetInvalidRepo()

	// Format data.
	output := Output{
		Time:    time.Now(),
		Data:    awesomeRepos,
		Invalid: invalidRepos,
	}
	outputBytes, err := json.Marshal(output)
	if err != nil {
		errMsg := "Failed to format data."
		fmt.Fprintln(writer, errMsg, strerr(err))
		logger.Error(errMsg, zap.Error(err))
		return err
	}

	// Output data.
	if len(w.outputPath) != 0 {
		err := ioutil.WriteFile(w.outputPath, outputBytes, 0644)
		if err != nil {
			errMsg := "Failed to output data."
			fmt.Fprintln(writer, errMsg, strerr(err))
			logger.Error(errMsg, zap.Error(err))
		}
	} else {
		fmt.Fprintln(writer, string(outputBytes))
	}

	// Warning message.
	if len(invalidRepos) > 0 {
		fmt.Fprintf(writer, "\nCatch some invalid repositories: %v\n", invalidRepos)
	}

	// The last message.
	if len(w.outputPath) > 0 {
		fmt.Fprintf(writer, "Done. Output file: %s\n", w.outputPath)
	} else {
		fmt.Fprintln(writer, "Done.")
	}
	return nil
}

func (w *Worker) newReporter() *awg.Reporter {
	return &awg.Reporter{}
}

func (w *Worker) newAwgClient(config config.Config) (*awg.Client, error) {
	logger := w.logger
	defer logger.Sync()

	htmlClient := cohttp.NewClient(http.Client{},
		config.MaxConcurrent, config.Network.RetryTime,
		config.Network.RetryInterval, config.LogRespHead, nil)
	apiClient := cohttp.NewClient(http.Client{},
		config.MaxConcurrent, config.Network.RetryTime,
		config.Network.RetryInterval, config.LogRespHead, w.reporter)

	options := github.NewDefaultClientOption()
	options.AccessToken = config.AccessToken
	if len(config.Github.HTMLHost) > 0 {
		options.HTMLHost = config.Github.HTMLHost
	}
	if len(config.Github.ApiHost) > 0 {
		options.APIHost = config.Github.ApiHost
	}

	gc, err := github.NewClient(htmlClient, apiClient, options)
	if err != nil {
		errMsg := "Failed to create github client."
		fmt.Fprintln(w.writer, errMsg)
		logger.Error(errMsg, zap.Error(err))
		return nil, err
	}
	client, err := awg.NewClient(gc)
	if err != nil {
		errMsg := "Failed to create awg client."
		fmt.Fprintln(w.writer, errMsg)
		logger.Error(errMsg, zap.Error(err))
		return nil, err
	}
	return client, nil
}

func (w *Worker) progressBar(prefix string) (completed <-chan struct{}, cancel chan<- struct{}) {
	pbCompleted := make(chan struct{})
	pbCancel := make(chan struct{})

	// TODO: refactor later
	getTotalNum := func() (numTotal int, canceled bool) {
		ticker := time.NewTicker(time.Second)
		for {
			select {
			case <-ticker.C:
				numTotal = w.reporter.GetTotalRepoNum()
				if numTotal > 0 {
					return numTotal, false
				}
			case <-pbCancel:
				return 0, true
			}
		}
	}
	go func() {
		numTotal, canceled := getTotalNum()
		if canceled {
			return
		}
		bar := progressbar.NewOptions(numTotal,
			progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
			progressbar.OptionEnableColorCodes(true),
			progressbar.OptionSetWidth(15),
			progressbar.OptionSetDescription(prefix),
			progressbar.OptionShowIts(),
			progressbar.OptionShowCount(),
		)
		ticker := time.NewTicker(time.Second)
		for {
			select {
			case <-pbCancel:
				close(pbCompleted)
				return
			case <-ticker.C:
				numCompleted := w.reporter.GetFinishedRepoNum()
				if numCompleted >= numTotal {
					bar.Finish()
					close(pbCompleted)
					return
				}
				bar.Set(numCompleted)
			}
		}
	}()
	return pbCompleted, pbCancel
}
