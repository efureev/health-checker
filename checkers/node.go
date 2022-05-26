package checkers

import (
	"fmt"
	"os/exec"

	"github.com/efureev/health-checker"
)

func CheckingNodeFn(min string) checker.CmdFn {
	return func(result *checker.Result, log checker.ILogger) {
		out, err := exec.Command("node", `-v`).Output()
		if err != nil {
			result.AddErrorFromStr(`Node is missing`)
			return
		}
		if string(out) < `v`+min {
			result.AddErrorFromStr(fmt.Sprintf(`Need a min version: %s. Found: %s`, min, out))
			return
		}

		result.Info.Version = string(out)

		return
	}
}
