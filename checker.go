package checker

import (
	"fmt"
	"sync"

	"github.com/efureev/go-multierror"
)

type Checker struct {
	list   []ICmd
	errors multierror.Collection
	lgr    ILogger
}

func (c *Checker) AddCmd(cmd ...ICmd) *Checker {
	c.list = append(c.list, cmd...)
	return c
}

func (c Checker) receiveErrors() error {
	if c.errors.HasErrors() {
		return c.errors
	}

	return nil
}

func (c *Checker) Run() error {
	for _, cmd := range c.list {
		cmd.SetLogger(c.lgr).Run()
		c.addError(cmd.Result().Error())
	}

	return c.receiveErrors()
}

func (c *Checker) RunParallel() error {
	wg := sync.WaitGroup{}

	for _, cmd := range c.list {
		if !cmd.Enable() {
			continue
		}

		wg.Add(1)

		go func(cmd ICmd, l ILogger) {
			defer wg.Done()

			l.Info(`[%s] checking...`, cmd.Name())
			cmd.SetLogger(l).Run()
			c.addError(multierror.Prefix(cmd.Result().Error(), cmd.Name()+`:`))
			printResult(cmd, l)
		}(cmd, c.lgr)
	}

	wg.Wait()

	return c.receiveErrors()
}

func (c *Checker) SetLogger(l ILogger) *Checker {
	c.lgr = l
	return c
}

func (c *Checker) addError(err error) {
	if err != nil {
		c.errors.Append(err)
	}
}

func NewChecker() *Checker {
	return &Checker{}
}

func printResult(cmd ICmd, l ILogger) {
	result := cmd.Result()

	msg := fmt.Sprintf(`[%s] %s`, cmd.Name(), result.Status)
	switch result.Status {
	case StatusDone:
		l.Success(msg)
		l.Log(result.Info.String())
	case StatusFailed:
		l.Error(msg)
		l.Log(result.Error().Error())
	case StatusPending:
		l.Log(msg)
	}
}
