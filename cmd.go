package checker

type CmdFn func(result *Result, log ILogger)

type ICmd interface {
	Enable() bool
	Name() string
	Run()
	Result() *Result
	SetLogger(l ILogger) *Cmd
}

type Cmd struct {
	enable bool
	name   string
	fn     CmdFn
	result *Result
	out    ILogger
}

func (cmd Cmd) Enable() bool {
	return cmd.enable
}

func (cmd Cmd) Result() *Result {
	return cmd.result
}

func (cmd *Cmd) Run() {
	if cmd.fn != nil {
		cmd.result.StatusPending()
		cmd.fn(cmd.result, cmd.out)

		if cmd.result.HasError() {
			cmd.result.StatusFailed()
		} else {
			cmd.result.StatusDone()
		}
	}
}

func (cmd *Cmd) SetCheckFn(fn CmdFn) *Cmd {
	cmd.fn = fn
	return cmd
}

func (cmd *Cmd) SetEnable(enable bool) *Cmd {
	cmd.enable = enable
	return cmd
}

func (cmd Cmd) Name() string {
	return cmd.name
}

func (cmd *Cmd) SetLogger(l ILogger) *Cmd {
	cmd.out = l
	return cmd
}

func NewCmd(name string) *Cmd {
	return &Cmd{name: name, enable: true, result: &Result{}}
}
