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

func IsSameVersion(builders []Builder) (bool, error) {
	sameVersion := true
	for _, b := range builders {
		vi, err := b.InstalledVersion()
		if err != nil {
			return false, err
		}
		switch b.Component {
		case ComponentPcre:
			fallthrough
		case ComponentOpenSSL:
			fallthrough
		case ComponentZlib:
			if vi == "" && !b.Static {
				continue
			} else if vi == b.Version && b.Static {
				continue
			}
		default:
			if vi == b.Version {
				continue
			}
		}
		sameVersion = false
	}

	if sameVersion {
		return true, nil
	}

	return false, nil
}
