package configure

import (
	"bufio"
	"os"

	"github.com/cubicdaiya/nginx-build/command"
)

func Run() error {
	args := []string{"sh", "./nginx-configure"}
	if command.VerboseEnabled {
		return command.Run(args)
	}

	f, err := os.Create("nginx-configure.log")
	if err != nil {
		return command.Run(args)
	}
	defer f.Close()

	cmd, err := command.Make(args)
	if err != nil {
		return err
	}

	writer := bufio.NewWriter(f)
	cmd.Stdout = writer
	defer writer.Flush()

	return cmd.Run()
}
