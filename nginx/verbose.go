package nginx

import (
	"os"
	"os/exec"
)

var Verboseenabled bool

func checkVerboseenabled(cmd *exec.Cmd) {
	if Verboseenabled {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
}
