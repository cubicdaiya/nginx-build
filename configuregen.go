package main

import (
	"fmt"
	"runtime"
	"strings"
)

type StaticLibrary struct {
	Name    string
	Version string
	Option  string
}

func makeStaticLibrary(builder *Builder) StaticLibrary {
	return StaticLibrary{
		Name:    builder.name(),
		Version: builder.Version,
		Option:  builder.option()}
}

func configureGenModule3rd(modules3rd []Module3rd) string {
	result := ""
	for _, m := range modules3rd {
		if m.Form == "local" {
			result += fmt.Sprintf("--add-module=%s \\\n", m.Url)
		} else {
			result += fmt.Sprintf("--add-module=../%s \\\n", m.Name)
		}
	}
	return result
}

func configureGen(configure string, modules3rd []Module3rd, dependencies []StaticLibrary) string {
	openSSLStatic := false
	if len(configure) == 0 {
		configure = `#!/bin/sh

./configure \
`
		if runtime.GOOS == "darwin" {
			configure += "--with-cc-opt=\"-Wno-deprecated-declarations\" \\"
		}
	}

	for _, d := range dependencies {
		configure += fmt.Sprintf("%s=../%s-%s \\\n", d.Option, d.Name, d.Version)
		if d.Name == "openssl" {
			openSSLStatic = true
		}
	}

	if openSSLStatic && !strings.Contains(configure, "--with-http_ssl_module") {
		configure += "--with-http_ssl_module \\\n"
	}

	configure_modules3rd := configureGenModule3rd(modules3rd)
	configure += configure_modules3rd

	return configure
}
