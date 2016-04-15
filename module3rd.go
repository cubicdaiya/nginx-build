package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-ini/ini"
)

type Module3rd struct {
	Name      string
	Form      string
	Url       string
	Rev       string
	Dynamic   bool
	Shprov    string
	ShprovDir string
}

func loadModule3rd(s *ini.Section) Module3rd {
	var (
		module3rd Module3rd
	)

	module3rd.Name = s.Name()
	module3rd.Form = s.Key("form").String()
	module3rd.Url = s.Key("url").String()
	module3rd.Rev = s.Key("rev").String()
	module3rd.Shprov = s.Key("shprov").String()
	module3rd.ShprovDir = s.Key("shprovdir").String()
	module3rd.Dynamic = s.Key("dynamic").MustBool()

	return module3rd
}

func loadModules3rd(f *ini.File) []Module3rd {
	var modules3rd []Module3rd
	for _, s := range f.Sections() {
		if s.Name() == "DEFAULT" {
			continue
		}
		modules3rd = append(modules3rd, loadModule3rd(s))
	}
	return modules3rd
}

func loadModules3rdFile(path string) ([]Module3rd, error) {
	var modules3rd []Module3rd
	if len(path) > 0 {
		if !fileExists(path) {
			return modules3rd, fmt.Errorf("modulesConfPath(%s) does not exist.", path)
		}
		modulesConf, err := ini.Load(path)
		if err != nil {
			return modules3rd, err
		}
		modules3rd = loadModules3rd(modulesConf)
	}
	return modules3rd, nil
}

func provideShell(sh string) error {
	if len(sh) == 0 {
		return nil
	}

	cmds := strings.Split(strings.Trim(sh, " "), "&&")

	for _, cmd := range cmds {
		args := strings.Split(strings.Trim(cmd, " "), " ")
		if err := runCommand(args); err != nil {
			return err
		}
	}

	return nil
}

func provideModule3rd(m *Module3rd) error {
	if len(m.Rev) > 0 {
		dir := saveCurrentDir()
		os.Chdir(m.Name)
		err := switchRev(m.Form, m.Rev)
		if err != nil {
			return fmt.Errorf("%s (%s checkout %s): %s", m.Name, m.Form, m.Rev, err.Error())
		}
		os.Chdir(dir)
	}

	if len(m.Shprov) > 0 {
		dir := saveCurrentDir()
		if len(m.ShprovDir) > 0 {
			os.Chdir(m.Name + "/" + m.ShprovDir)
		} else {
			os.Chdir(m.Name)
		}
		err := provideShell(m.Shprov)
		if err != nil {
			return fmt.Errorf("%s's shprov(%s): %s", m.Name, m.Shprov, err.Error())
		}
		os.Chdir(dir)
	}

	return nil
}

func switchRev(form, rev string) error {
	var err error

	switch form {
	case "git":
		err = runCommand([]string{"git", "checkout", rev})
	case "hg":
		err = runCommand([]string{"hg", "checkout", rev})
	default:
		err = fmt.Errorf("form=%s is not supported", form)
	}

	return err
}
