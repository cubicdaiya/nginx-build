package builder

import (
	"bufio"
	"os"
	"strconv"

	"github.com/cubicdaiya/nginx-build/command"
)

func BuildNginx(jobs int) error {
	args := []string{"make", "-j", strconv.Itoa(jobs)}
	if command.VerboseEnabled {
		return command.Run(args)
	}

	f, err := os.Create("nginx-build.log")
	if err != nil {
		return command.Run(args)
	}
	defer f.Close()

	cmd, err := command.Make(args)
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(f)
	cmd.Stderr = writer
	defer writer.Flush()

	return cmd.Run()
}
