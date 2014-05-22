package common

import (
	"os"
	"os/exec"
)

var Verboseenabled bool

func CheckVerboseEnabled(cmd *exec.Cmd) {
	if Verboseenabled {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
}
