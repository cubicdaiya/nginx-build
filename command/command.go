package command

import (
	"errors"
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

func Make(args []string) (*exec.Cmd, error) {
	var cmd *exec.Cmd
	switch len(args) {
	case 0:
		return nil, errors.New("empty command")
	case 1:
		cmd = exec.Command(args[0])
	default:
		cmd = exec.Command(args[0], args[1:]...)
	}

	return cmd, nil
}

func Run(args []string) error {
	cmd, err := Make(args)
	if err != nil {
		return err
	}

	checkVerboseEnabled(cmd)
	return cmd.Run()
}
