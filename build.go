package main

import (
	_ "fmt"
	"os/exec"
	"strconv"
	"strings"
)

func configure() error {
	return runCommand(exec.Command("sh", "./nginx-configure"))
}

func build(jobs int) error {
	return runCommand(exec.Command("make", "-j", strconv.Itoa(jobs)))
}

func extractArchive(path string) error {
	return runCommand(exec.Command("tar", "zxvf", path))
}

func switchRev(rev string) error {
	return runCommand(exec.Command("git", "checkout", rev))
}

func provideShell(sh string) error {
	if len(sh) == 0 {
		return nil
	}
	args := strings.Split(strings.Trim(sh, " "), " ")
	var err error
	if len(args) == 1 {
		err = runCommand(exec.Command(args[0]))
	} else {
		err = runCommand(exec.Command(args[0], args[1:]...))
	}
	return err
}
