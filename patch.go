package main

import (
	"fmt"
	"os"

	"github.com/cubicdaiya/nginx-build/command"
)

func patch(path, option string) error {
	args := []string{"sh", "-c"}
	args = append(args, fmt.Sprintf("patch %s < %s", option, path))

	cmd, err := command.Make(args)
	if err != nil {
		return err
	}

	// As the output of patch is interactive,
	// the result is always output.
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
