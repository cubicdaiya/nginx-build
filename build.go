package main

import (
	_ "fmt"
	"os/exec"
	"strconv"
	"strings"
)

func configure() error {
	cmd := exec.Command("sh", "./nginx-configure")
	checkVerboseEnabled(cmd)
	return cmd.Run()
}

func make(jobs int) error {
	cmd := exec.Command("make", "-j", strconv.Itoa(jobs))
	checkVerboseEnabled(cmd)
	return cmd.Run()
}

func extractArchive(path string) error {
	cmd := exec.Command("tar", "zxvf", path)
	checkVerboseEnabled(cmd)
	return cmd.Run()
}

func switchRev(rev string) error {
	cmd := exec.Command("git", "co", rev)
	checkVerboseEnabled(cmd)
	return cmd.Run()
}

func prevShell(sh string) error {
	if len(sh) == 0 {
		return nil
	}
	args := strings.Split(strings.Trim(sh, " "), " ")
	var cmd *exec.Cmd
	if len(args) == 1 {
		cmd = exec.Command(args[0])
	} else {
		cmd = exec.Command(args[0], args[1:]...)
	}
	checkVerboseEnabled(cmd)
	return cmd.Run()
}
