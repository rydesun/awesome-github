你好！
很冒昧用这样的方式来和你沟通，如有打扰请忽略我的提交哈。我是光年实验室（gnlab.com）的HR，在招Golang开发工程师，我们是一个技术型团队，技术氛围非常好。全职和兼职都可以，不过最好是全职，工作地点杭州。
我们公司是做流量增长的，Golang负责开发SAAS平台的应用，我们做的很多应用是全新的，工作非常有挑战也很有意思，是国内很多大厂的顾问。
如果有兴趣的话加我微信：13515810775  ，也可以访问 https://gnlab.com/，联系客服转发给HR。
# Awesome GitHub

[![Build](https://github.com/rydesun/awesome-github/workflows/Build/badge.svg)](https://github.com/rydesun/awesome-github/actions?query=workflow%3ABuild)
[![Unit-Tests](https://github.com/rydesun/awesome-github/workflows/Unit-Tests/badge.svg)](https://github.com/rydesun/awesome-github/actions?query=workflow%3AUnit-Tests)
[![Coverage Status](https://coveralls.io/repos/github/rydesun/awesome-github/badge.svg?branch=master)](https://coveralls.io/github/rydesun/awesome-github?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/rydesun/awesome-github)](https://goreportcard.com/report/github.com/rydesun/awesome-github)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/rydesun/awesome-github/blob/master/LICENSE)

Explore your favorite Awesome GitHub repositories via the command-line tool awg!

[\[ Chinese | 中文 \]](https://github.com/rydesun/awesome-github/blob/master/README-zh_CN.md)

## Awesome Lists

What is Awesome Lists？

Such as famous [awesome-go](https://github.com/avelino/awesome-go), which is a member
of Awesome Lists, where we can quickly find frameworks, libraries, software,
and other resources related to Go.

Furthermore, we can find more awesome things from
[Awesome Lists](https://github.com/topics/awesome).

## Command-line Tool awg

At the moment, Awesome List usually exhibit lots of GitHub repositories.
For example, awesome-go contains thousands of GitHub repositories.
However, an Awesome List doesn't include information like the count of stars
or the last commit date for those repositories.
In many cases, we demand that information and have to manually open links
to see the them.

The command-line tool, awg, helps us dig a little deeper into the Awesome List
to find out more about all the GitHub repositories on it!

awg will fetch information of GitHub repositories listed on an Awesome List,
and output the data to a file of your choice.
You can use awg to generate a [browser page](#view-in-browser) to view the data later,
or use your favorite tools like jq or python to analyze.

![Screenshot](https://user-images.githubusercontent.com/19602440/88459895-f3897480-ce87-11ea-8fe7-13773037c56d.gif)

The final output:

```javascript
{
  "data": {
    "Command Line": [
      {
        "id": {
          "owner": "urfave",
          "name": "cli"
        },
        "owner": "urfave",
        "awesome_name": "urfave/cli",
        "link": "https://github.com/urfave/cli",
        "watch": 295,
        "star": 14171,
        "fork": 1134,
        "last_commit": "2020-07-12T13:32:01Z",
        "description": "A simple, fast, and fun package for building command line apps in Go",
        "awesome_description": "urfave/cli - Simple, fast, and fun package for building command line apps in Go (formerly codegangsta/cli)."
      },
      // ...
    ]
    // ...
  }
}
```

### Data Analysis

We can use any tool to analyze the obtained data file.

For example, after getting data file `awg.json`, which is concerned with awesome-go

#### View in Browser

The page effect is shown below

![Screenshot](https://user-images.githubusercontent.com/19602440/89290996-3fd37200-d649-11ea-8807-a6a117d016f0.png)

By runing this command

```bash

awg view awg.json
```

It will be running a simple web server locally,
listening at `127.0.0.1:3000` by default.
Open this page with your browser to view.
Can use `--listen` to specify other address.

Note: This does not mean that it can be viewed offline.
An Internet connection is necessary.

To replace the implementation for viewing,
you can replace the embedded JS script address via `--script`,
which supports local path.

#### View in Terminal

By using the popular command line tool jq, we can:

View [Command Line section](https://github.com/avelino/awesome-go#command-line)
of awesome-go, sorted by the count of stars

```bash
cat awg.json | jq '.data | ."Command Line" | sort_by(.star)'
```

## Installation

Get command-line tool awg

```bash
go get github.com/rydesun/awesome-github/cmd/awg
```

## Fetch Data

First prepare the following before running awg:

- GitHub personal access token
- awg configuration file

### Access Token

awg fetch information of GitHub repository by calling the GitHub GraphQL API.
This official API requires your personal access token to be verified
before you can use it.
So, you need to provide awg with a GitHub personal access token.

If you do not have the token, please view the article
[Creating a personal access token](https://docs.github.com/en/github/authenticating-to-github/creating-a-personal-access-token)。

awg does not need the token to have any scopes or permissions.
Therefore, do not grant the token any scopes or permissions.

### Configuration

A configuration file is necessary.
See [configuration file template](https://github.com/rydesun/awesome-github/blob/master/configs/config.yaml)
in the directory `configs` as a reference.

awg will read the personal access token from the
environment variable `GITHUB_ACCESS_TOKEN` first.
So it is not necessary to store this value in the configuration file.

Increasing the number of concurrent query requests will
increase the speed of the query,
but the number should not be too large (the current recommended value is 3),
otherwise this behavior will be viewed as abusing, and then blocked.

All relative paths in a configuration file are relative to
awg's current working directory,
not to the directory where the configuration file is located.
Absolute paths are encouraged.

### Run

Fetch JSON data

```bash
awg fetch --config path/to/config.yaml
```

(Recommended) Specify the GitHub Personal Access Token from an environment variable

```bash
GITHUB_ACCESS_TOKEN=<Your Token> awg fetch --config path/to/config.yaml
```

Note the rate limit, which is not identical to the number of concurrent requests.
Currently awg can query up to 5000 GitHub repositories per hour.
If there are too many queries, GitHub will impose limits
on requests and responses.
See [GitHub Resource limitations](https://docs.github.com/en/graphql/overview/resource-limitations#rate-limit) for more.

## Notes

This project is currently only testing awesome-go's lists,
other testing for Awesome List results are pending.

awg does not support Windows platforms.
