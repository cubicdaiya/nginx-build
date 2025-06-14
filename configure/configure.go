package configure

import (
	"bufio"
	"fmt"
	"log"
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
		log.Printf("[warn] could not create nginx-configure.log: %v", err)
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

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("configure failed: %w", err)
	}

	return nil
}
