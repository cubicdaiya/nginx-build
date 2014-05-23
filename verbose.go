package main

import (
	"os"
	"os/exec"
)

var VerboseEnabled bool

func checkVerboseEnabled(cmd *exec.Cmd) {
	if VerboseEnabled {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
}
