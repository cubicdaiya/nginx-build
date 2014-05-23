package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type StaticLibrary struct {
	Name    string
	Version string
	Option  string
}

func configureGenModule3rd(modules3rd []Module3rd) string {
	result := ""
	//for i := 0; i < len(modules3rd); i++ {
	for _, m := range modules3rd {
		result += fmt.Sprintf("--add-module=../%s \\\n", m.Name)
	}
	return result
}

func (builder *Builder) configureGen(conf string, modules3rd []Module3rd, dependencies []StaticLibrary) error {
	configure := `#!/bin/sh

./configure `

	if conf != "" {
		configure += "\\\n"
		options := strings.Split(conf, "\n")

		for i := range options {
			options[i] += " \\"
		}

		conf = strings.Join(options, "\n")
		configure += conf
	}

	for _, d := range dependencies {
		configure += fmt.Sprintf("%s=../%s-%s \\\n", d.Option, d.Name, d.Version)
	}

	configure_modules3rd := configureGenModule3rd(modules3rd)
	configure += configure_modules3rd

	return ioutil.WriteFile("./nginx-configure", []byte(configure), 0655)
}
