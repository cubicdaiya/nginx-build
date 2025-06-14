package module3rd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/cubicdaiya/nginx-build/command"
	"github.com/cubicdaiya/nginx-build/util"
)

func Provide(m *Module3rd) error {
	if len(m.Rev) > 0 {
		dir := util.SaveCurrentDir()
		if err := os.Chdir(m.Name); err != nil {
			return fmt.Errorf("chdir to %s failed: %w", m.Name, err)
		}
		if err := switchRev(m.Form, m.Rev); err != nil {
			return fmt.Errorf("%s (%s checkout %s): %s", m.Name, m.Form, m.Rev, err.Error())
		}
		if err := os.Chdir(dir); err != nil {
			return fmt.Errorf("return to dir %s failed: %w", dir, err)
		}
	}

	if len(m.Shprov) > 0 {
		dir := util.SaveCurrentDir()
		if len(m.ShprovDir) > 0 {
			if err := os.Chdir(m.Name + "/" + m.ShprovDir); err != nil {
				return fmt.Errorf("chdir to %s/%s failed: %w", m.Name, m.ShprovDir, err)
			}
		} else {
			if err := os.Chdir(m.Name); err != nil {
				return fmt.Errorf("chdir to %s failed: %w", m.Name, err)
			}
		}
		if err := provideShell(m.Shprov); err != nil {
			return fmt.Errorf("%s's shprov(%s): %s", m.Name, m.Shprov, err.Error())
		}
		if err := os.Chdir(dir); err != nil {
			return fmt.Errorf("return to dir %s failed: %w", dir, err)
		}
	}

	return nil
}

func provideShell(sh string) error {
	if strings.TrimSpace(sh) == "" {
		return nil
	}
	if command.VerboseEnabled {
		return command.Run([]string{"sh", "-c", sh})
	}

	cmd := exec.Command("sh", "-c", sh)
	return cmd.Run()
}

func switchRev(form, rev string) error {
	var err error

	switch form {
	case "git":
		err = command.Run([]string{"git", "checkout", rev})
	case "hg":
		err = command.Run([]string{"hg", "checkout", rev})
	default:
		err = fmt.Errorf("form=%s is not supported", form)
	}

	return err
}
