# Health checker

[![Go Report Card](https://goreportcard.com/badge/github.com/efureev/health-checker)](https://goreportcard.com/report/github.com/efureev/health-checker)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/efureev/health-checker)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/efureev/health-checker)
![GitHub](https://img.shields.io/github/license/efureev/health-checker)

## Installation

```shell
go get github.com/efureev/health-checker
```

## Description

Health checker allows you check out various services and your custom processes for your custom checks (i.e.
availability).

## Usage

Create checker:

```go
ch := checker.NewChecker().
SetLogger(logger.NewConsoleLogger()).
```

Add a command to healthcheck:

```go
package main

import (
	checker "github.com/efureev/health-checker"
	"github.com/efureev/health-checker/checkers"
)

func main() {
	ch := checker.NewChecker().
		// SetLogger(logger.NewConsoleLogger()).
		AddCmd(
			checker.NewCmd(`Node`).SetCheckFn(checkers.CheckingNodeFn(`16`)),
			checker.NewCmd(`Redis`).SetCheckFn(
				checkers.CheckingRedisFn(checkers.RedisConfig{
					Host:     `localhost`,
					Password: ``,
					DB:       0,
					Port:     6379,
				}),
			),
		)

	// aSynchro
	err := ch.RunParallel()

	// synchro
	err := ch.Run()

}
```

With [`Commander`](https://github.com/efureev/go-commander):

```go
package main

import (
	commander "github.com/efureev/go-commander"
	checker "github.com/efureev/health-checker"
	"github.com/efureev/health-checker/checkers"
)

func main() {
	cnf := struct{}{}
	cmdr := commander.NewCommander().
		Add(
			cmdChecking(cnf),
			// ... 
		)

	cmdr.Run()

	if cmdr.HasError() {
		utils.Error(cmdr.commander.Error(), "\nErrors:\n")
	}
}

func cmdChecking(cnf any) *commander.Command {
	return commander.NewCommand(`checking`).
		OnRun(func(cmd *commander.Command) error {
			ch := checker.NewChecker().
				SetLogger(logger.NewConsoleLogger()).
				AddCmd(
					checkingPostgresDb(cnf.Checker.Database),
					checkingElasticsearch(cnf.Checker.Elasticsearch),
					checkingRedis(cnf.Checker.Redis),
					checkingRabbitMQ(cnf.Checker.Rabbit),
					checkingCaptcha(cnf.Checker.Captcha),
					checkingTargetRepos(cnf.Targets),
					checkingHDD(cnf.Folder),
					checkingNode(),
				)
			return ch.RunParallel()
		})
}
```
