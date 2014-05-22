package main

import (
        "github.com/robfig/config"
)

type Module3rd struct {
        Name string
        Form string
        Url  string
        Rev  string
}

func LoadModule3rd(name string, c config.Config) Module3rd {
        var module3rd Module3rd
        module3rd.Name = name
        module3rd.Form, _ = c.String(name, "form")
        module3rd.Url, _ = c.String(name, "url")
        module3rd.Rev, _ = c.String(name, "rev")
        return module3rd
}

func LoadModules3rd(c *config.Config) []Module3rd {
        sections := c.Sections()
        l := len(sections)
        var modules3rd []Module3rd
        for i := 0; i < l; i++ {
                if sections[i] == config.DEFAULT_SECTION {
                        continue
                }
                modules3rd = append(modules3rd, LoadModule3rd(sections[i], *c))
        }
        return modules3rd
}
