# Awesome GitHub

[![Build](https://github.com/rydesun/awesome-github/workflows/Build/badge.svg)](https://github.com/rydesun/awesome-github/actions?query=workflow%3ABuild)
[![Unit-Tests](https://github.com/rydesun/awesome-github/workflows/Unit-Tests/badge.svg)](https://github.com/rydesun/awesome-github/actions?query=workflow%3AUnit-Tests)
[![Coverage Status](https://coveralls.io/repos/github/rydesun/awesome-github/badge.svg?branch=master)](https://coveralls.io/github/rydesun/awesome-github?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/rydesun/awesome-github)](https://goreportcard.com/report/github.com/rydesun/awesome-github)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/rydesun/awesome-github/blob/master/LICENSE)

通过命令行工具 awg，探索你钟爱的 Awesome GitHub 仓库！

## Awesome Lists

什么是 Awesome Lists？

比如著名的 [awesome-go](https://github.com/avelino/awesome-go) 就是 Awesome Lists 的一员，
我们可以从中快速找到很多和 Go 相关的框架、库、软件等资源。

不止是 Go，还可以从 [Awesome Lists](https://github.com/topics/awesome)
寻找更多你感兴趣的内容。

## 命令行工具 awg

当前，Awesome List 一般会列举很多 GitHub 仓库，
就比如 awesome-go 包含了上千个 GitHub 仓库。
但是，Awesome List 不包含这些仓库的 star 数、
更新时间 (最后一次 commit 时间) 之类的信息。
很多时候，我们需要这些信息作为参考，只能手动打开链接以查看这些仓库的信息。

命令行工具 awg，帮助我们进一步挖掘 Awesome List 中所有 GitHub 仓库的惊人信息！

awg 将一次性获取指定 Awesome List 中的 GitHub 仓库信息，
并输出获取到的数据到指定文件中。
稍后你可以使用喜欢的数据处理工具，比如 jq 和 python，对这些数据进行分析。

![Screenshot](https://user-images.githubusercontent.com/19602440/88459895-f3897480-ce87-11ea-8fe7-13773037c56d.gif)

最终输出的文件内容：

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
        "star": 14171,
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

## 安装

获取命令行工具 awg

```bash
go get github.com/rydesun/awesome-github/cmd/awg
```

## 使用

在运行 awg 之前，先准备好：

- GitHub personal access token
- awg 配置文件

### Access Token

awg 通过调用 GitHub GraphQL API 获取 GitHub 仓库信息，
该官方 API 需要验证你的 personal access token 后才能使用。
所以，需要向 awg 提供一个 GitHub personal access token。

如果没有该 token，请先
[创建 personal access token](https://docs.github.com/en/github/authenticating-to-github/creating-a-personal-access-token)。

注意！awg 不需要该 token 有任何作用域或权限。所以，不要授予该 token 任何作用域或权限。

### 配置

必须准备一个配置文件。可以参考目录`configs`中的
[配置文件模板](https://github.com/rydesun/awesome-github/blob/master/configs/config.yaml)。

awg 会优先从环境变量`GITHUB_ACCESS_TOKEN`中读取 personal access token，
所以可以不用将该值储存在配置文件中。

提升并发查询请求数可以提升查询速度，但数量不要过大 (当前推荐值为 3)，
否则会被 GitHub 视作滥用 API 的行为而遭到临时封禁。

### 运行

获取 JSON 数据文件

```bash
awg fetch --config path/to/config.yaml
```

(推荐) 从环境变量中指定 GitHub Personal Access Token 的形式运行

```bash
GITHUB_ACCESS_TOKEN=<Your Token> awg fetch --config path/to/config.yaml
```

请注意速率限制 (RateLimit)，该值不是并发请求数。
当前 awg 每小时最多查询 5000 个 GitHub 仓库。
如果查询次数过多，会受到 GitHub 的限制从而导致失败。
具体信息请参考 [GitHub Resource limitations](https://docs.github.com/en/graphql/overview/resource-limitations#rate-limit)。

## 数据分析

可以使用任何工具去分析获得的数据文件。

比如在获取 awesome-go 的数据文件`awg.json`后，
通过使用流行的命令行工具 jq，你可以：

查看 awesome-go 列表中 [Command Line](https://github.com/avelino/awesome-go#command-line)
一节的内容，并按照仓库的 star 数进行排序

```bash
cat awg.json | jq '.data | ."Command Line" | sort_by(.star)'
```

## 注意事项

当前该项目仅测试了 awesome-go 的列表，其他 Awesome List 的结果待检验。

awg 不支持 Windows 平台。
