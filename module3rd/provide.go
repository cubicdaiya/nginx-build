package module3rd

import (
	"fmt"
	"os"
	"strings"

	"github.com/cubicdaiya/nginx-build/command"
	"github.com/cubicdaiya/nginx-build/util"
)

func Provide(m *Module3rd) error {
	if len(m.Rev) > 0 {
		dir := util.SaveCurrentDir()
		os.Chdir(m.Name)
		if err := switchRev(m.Form, m.Rev); err != nil {
			return fmt.Errorf("%s (%s checkout %s): %s", m.Name, m.Form, m.Rev, err.Error())
		}
		os.Chdir(dir)
	}

	if len(m.Shprov) > 0 {
		dir := util.SaveCurrentDir()
		if len(m.ShprovDir) > 0 {
			os.Chdir(m.Name + "/" + m.ShprovDir)
		} else {
			os.Chdir(m.Name)
		}
		if err := provideShell(m.Shprov); err != nil {
			return fmt.Errorf("%s's shprov(%s): %s", m.Name, m.Shprov, err.Error())
		}
		os.Chdir(dir)
	}

	return nil
}

func provideShell(sh string) error {
	if len(sh) == 0 {
		return nil
	}

	cmds := strings.Split(strings.Trim(sh, " "), "&&")

	for _, cmd := range cmds {
		args := strings.Split(strings.Trim(cmd, " "), " ")
		if err := command.Run(args); err != nil {
			return err
		}
	}

	return nil
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
