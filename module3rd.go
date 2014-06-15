package main

import (
	"fmt"
	"github.com/robfig/config"
	"log"
	"os"
)

type Module3rd struct {
	Name   string
	Form   string
	Url    string
	Rev    string
	Shprov string
}

func loadModule3rd(name string, c *config.Config) Module3rd {
	var module3rd Module3rd
	module3rd.Name = name
	module3rd.Form, _ = c.String(name, "form")
	module3rd.Url, _ = c.String(name, "url")
	module3rd.Rev, _ = c.String(name, "rev")
	module3rd.Shprov, _ = c.String(name, "shprov")
	return module3rd
}

func loadModules3rd(c *config.Config) []Module3rd {
	sections := c.Sections()
	var modules3rd []Module3rd
	for _, s := range sections {
		if s == config.DEFAULT_SECTION {
			continue
		}
		modules3rd = append(modules3rd, loadModule3rd(s, c))
	}
	return modules3rd
}

func loadModules3rdFile(path string) ([]Module3rd, error) {
	var modules3rd []Module3rd
	if len(path) > 0 {
		if !fileExists(path) {
			return modules3rd, fmt.Errorf("modulesConfPath(%s) does not exist.", path)
		}
		modulesConf, err := config.ReadDefault(path)
		if err != nil {
			return modules3rd, err
		}
		modules3rd = loadModules3rd(modulesConf)
	}
	return modules3rd, nil
}

func provideModule3rd(m *Module3rd) {
	if len(m.Rev) > 0 {
		dir := saveCurrentDir()
		os.Chdir(m.Name)
		err := switchRev(m.Rev)
		if err != nil {
			log.Println(err.Error())
		}
		os.Chdir(dir)
	}

	if len(m.Shprov) > 0 {
		dir := saveCurrentDir()
		os.Chdir(m.Name)
		err := provideShell(m.Shprov)
		if err != nil {
			log.Println(err.Error())
		}
		os.Chdir(dir)
	}
}
