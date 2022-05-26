package checker

import (
	"fmt"
)

type Status int

const (
	StatusUnknown Status = iota
	StatusPending
	StatusDone
	StatusFailed
)

func (s Status) String() string {
	switch s {
	case StatusDone:
		return `Done`
	case StatusFailed:
		return `Failed`
	case StatusPending:
		return `Pending`
	default:
		return `Unknown`
	}
}

type Info struct {
	Version string
	Text    string
}

func (i Info) String() string {
	var msg string
	if i.Version != `` {
		msg += fmt.Sprintf("   Version: %s\n", i.Version)
	}

	if i.Text != `` {
		msg += fmt.Sprintf("   : %s\n", i.Text)
	}

	return msg
}

type resultStep struct {
	Text   string
	status Status
	error  error
}

func (s *resultStep) AddError(err error) {
	s.error = err
}

func (s *resultStep) Failed(err ...error) {
	s.status = StatusFailed
	if len(err) > 0 {
		s.AddError(err[0])
	}
}

func (s *resultStep) Success() {
	s.status = StatusDone
}

func newStep(t string) *resultStep {
	return &resultStep{Text: t, status: StatusUnknown}
}

type Result struct {
	Success string
	error   error
	steps   []*resultStep
	Status  Status
	Info    Info
}

func (r *Result) AddError(err error) *Result {
	r.error = err

	return r
}

func (r *Result) AddErrorFromStr(err string) *Result {
	return r.AddError(fmt.Errorf(err))
}

func (r *Result) StatusPending() {
	r.Status = StatusPending
}

func (r *Result) StatusFailed() {
	r.Status = StatusFailed
}

func (r *Result) StatusDone() {
	r.Status = StatusDone
}

func (r *Result) AddStep(step string) *resultStep {
	s := newStep(step)
	r.steps = append(r.steps, newStep(step))

	return s
}

func (r Result) HasError() bool {
	return r.error != nil
}

func (r Result) Error() error {
	return r.error
}
