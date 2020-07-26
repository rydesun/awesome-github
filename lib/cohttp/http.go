package cohttp

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/rydesun/awesome-github/lib/errcode"
)

type Reporter interface {
	ConReqNum(int)
}

type Client struct {
	c             *http.Client
	queue         chan struct{}
	MaxConcurrent int
	RetryTime     int
	RetryInterval time.Duration
	logRespHead   int
	reporter      Reporter
}

func NewClient(client http.Client, maxConcurrent int,
	retryTime int, retryInterval time.Duration,
	logRespHead int, reporter Reporter) *Client {
	var queue chan struct{}
	if maxConcurrent > 0 {
		queue = make(chan struct{}, maxConcurrent)
	}
	if retryTime < 0 {
		retryTime = 0
	}
	return &Client{
		c:             &client,
		queue:         queue,
		MaxConcurrent: maxConcurrent,
		RetryTime:     retryTime,
		RetryInterval: retryInterval,
		logRespHead:   logRespHead,
		reporter:      reporter,
	}
}

func (c *Client) Do(req *http.Request) (resp *http.Response, err error) {
	const funcIntent = "try to send http request"
	const funcErrMsg = "failed to send http request"
	var (
		method = req.Method
		url    = req.URL.String()
		logger = getLogger()
	)
	defer logger.Sync()
	logger.Debug(funcIntent,
		zap.String("method", method),
		zap.String("url", url))

	if c.queue == nil {
		if c.reporter != nil {
			c.reporter.ConReqNum(1)
		}
	} else {
		c.queue <- struct{}{}
		if c.reporter != nil {
			c.reporter.ConReqNum(len(c.queue))
		}
	}
	logger.Debug("request is being sent",
		zap.String("method", method),
		zap.String("url", url))
	resp, err = c.c.Do(req)
	if c.queue != nil {
		<-c.queue
	}
	if err != nil {
		logger.Error(funcErrMsg, zap.Error(err),
			zap.String("method", method),
			zap.String("url", url))
		err = errcode.New(funcErrMsg, ErrCodeNetwork, ErrScope,
			[]string{err.Error()})
		return
	}

	logger.Debug("receive a response",
		zap.String("method", method),
		zap.Int("status", resp.StatusCode),
		zap.String("url", url))
	return
}

func (c *Client) DoBetter(req *http.Request) (
	rawdata []byte, hasBody bool, err error) {
	const funcIntent = "try to send http request"
	const funcErrMsg = "failed to " + funcIntent
	var (
		method = req.Method
		url    = req.URL.String()
		logger = getLogger()
	)
	defer logger.Sync()
	logger.Debug(funcIntent,
		zap.String("method", method),
		zap.String("url", url))

	var resp *http.Response
	for i := 0; i < c.RetryTime+1; i++ {
		resp, err = c.Do(req)
		if err == nil {
			break
		}
		logger.Warn(funcErrMsg, zap.Error(err),
			zap.String("method", method),
			zap.String("url", url),
			zap.Int("retry", i))
		time.Sleep(c.RetryInterval)
	}
	if err != nil {
		logger.Error(funcErrMsg, zap.Error(err),
			zap.String("method", method),
			zap.String("url", url))
		err = errcode.Wrap(err, funcErrMsg)
		return
	}
	rawdata, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		errMsg := "failed to read response"
		logger.Error(errMsg, zap.Error(err))
		err = errcode.Wrap(err, errMsg)
		return
	}

	// Read body successfully
	hasBody = true

	statusCode := resp.StatusCode
	if statusCode < 200 || statusCode > 299 {
		errMsg := "remote server did not return ok"
		logger.Error(errMsg,
			zap.String("method", method),
			zap.String("url", url),
			zap.ByteString("recv", rawdata),
			zap.Int("status", statusCode))
		err = errcode.New(errMsg, errcode.ErrCode(statusCode),
			ErrScope, nil)
		return
	}
	return
}

func (c *Client) Text(req *http.Request) (string, error) {
	const funcIntent = "try to get text from remote server"
	const funcErrMsg = "failed to get text from remote server"
	var (
		method = req.Method
		url    = req.URL.String()
		logger = getLogger()
	)
	defer logger.Sync()
	logger.Debug(funcIntent,
		zap.String("method", method),
		zap.String("url", url))

	rawdata, hasBody, err := c.DoBetter(req)
	var text string
	if hasBody {
		text = string(rawdata)
	}
	if err != nil {
		logger.Error(funcErrMsg, zap.Error(err),
			zap.String("method", method),
			zap.String("url", url))
		err = errcode.Wrap(err, funcErrMsg)
		return text, err
	}
	return text, nil
}

func (c *Client) Json(req *http.Request, respJson interface{}) error {
	const funcIntent = "try to get json from remote server"
	const funcErrMsg = "failed to get json from remote server"
	var (
		method = req.Method
		url    = req.URL.String()
		logger = getLogger()
	)
	defer logger.Sync()
	logger.Debug(funcIntent,
		zap.String("method", method),
		zap.String("url", url))

	rawdata, hasBody, err := c.DoBetter(req)

	// impossible: hasBody == false && err == nil
	if !hasBody && err != nil {
		logger.Error(funcErrMsg, zap.Error(err),
			zap.String("method", method),
			zap.String("url", url))
		err = errcode.Wrap(err, funcErrMsg)
		return err
	}
	logger.Debug("receive rawdata from remote",
		zap.ByteString("content", rawdata))
	// DO NOT cover err
	_err := json.Unmarshal(rawdata, &respJson)
	if _err != nil {
		errMsg := "failed to parse response"
		length := len(rawdata)
		logField := []zap.Field{
			zap.Error(err),
			zap.Int("length", length),
			zap.String("method", method),
			zap.String("url", url),
		}
		if c.logRespHead > 0 {
			content := truncate(rawdata, c.logRespHead)
			logField = append(logField,
				zap.ByteString("content", content))
		}
		//
		if err == nil {
			err = errcode.New(errMsg, ErrCodeJson, ErrScope, nil)
		}
		logger.Error(errMsg, logField...)
		return err
	}
	// Must return err!
	return err
}
