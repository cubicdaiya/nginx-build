package main

import (
	"fmt"
	"github.com/robfig/config"
)

type Module3rd struct {
	Name string
	Form string
	Url  string
	Rev  string
	PrevSh string
}

func loadModule3rd(name string, c config.Config) Module3rd {
	var module3rd Module3rd
	module3rd.Name = name
	module3rd.Form, _ = c.String(name, "form")
	module3rd.Url, _ = c.String(name, "url")
	module3rd.Rev, _ = c.String(name, "rev")
	module3rd.PrevSh, _ = c.String(name, "prevsh")
	return module3rd
}

func loadModules3rd(c *config.Config) []Module3rd {
	sections := c.Sections()
	l := len(sections)
	var modules3rd []Module3rd
	for i := 0; i < l; i++ {
		if sections[i] == config.DEFAULT_SECTION {
			continue
		}
		modules3rd = append(modules3rd, loadModule3rd(sections[i], *c))
	}
	return modules3rd
}

func loadModules3rdFile(path string) ([]Module3rd, error) {
	var modules3rd []Module3rd
	if path != "" {
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
