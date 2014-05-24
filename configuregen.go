package main

import (
	"fmt"
	"io/ioutil"
)

type StaticLibrary struct {
	Name    string
	Version string
	Option  string
}

func configureGenModule3rd(modules3rd []Module3rd) string {
	result := ""
	for _, m := range modules3rd {
		result += fmt.Sprintf("--add-module=../%s \\\n", m.Name)
	}
	return result
}

func (builder *Builder) configureGen(configure string, modules3rd []Module3rd, dependencies []StaticLibrary) error {
	if len(configure) == 0 {
		configure = `#!/bin/sh

./configure \
`
	}

	for _, d := range dependencies {
		configure += fmt.Sprintf("%s=../%s-%s \\\n", d.Option, d.Name, d.Version)
	}

	configure_modules3rd := configureGenModule3rd(modules3rd)
	configure += configure_modules3rd

	return ioutil.WriteFile("./nginx-configure", []byte(configure), 0655)
}
