package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/cubicdaiya/nginx-build/builder"
)

func configureGenModule3rd(modules3rd []Module3rd) string {
	result := ""
	for _, m := range modules3rd {
		opt := "--add-module"
		if m.Dynamic {
			opt = "--add-dynamic-module"
		}
		if m.Form == "local" {
			result += fmt.Sprintf("%s=%s \\\n", opt, m.Url)
		} else {
			result += fmt.Sprintf("%s=../%s \\\n", opt, m.Name)
		}
	}
	return result
}

func configureGen(configure string, modules3rd []Module3rd, dependencies []builder.StaticLibrary, options ConfigureOptions, rootDir string) string {
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

	for _, option := range options.Values {
		if *option.Value != "" {
			if option.Name == "--add-module" {
				configure += normalizeAddModulePaths(*option.Value, rootDir, false)
			} else if option.Name == "--add-dynamic-module" {
				configure += normalizeAddModulePaths(*option.Value, rootDir, true)
			} else {
				if strings.Contains(*option.Value, " ") {
					configure += option.Name + "=" + "'" + *option.Value + "'" + " \\\n"
				} else {
					configure += option.Name + "=" + *option.Value + " \\\n"
				}
			}
		}
	}

	for _, option := range options.Bools {
		if *option.Enabled {
			configure += option.Name + " \\\n"
		}
	}

	return configure
}

func normalizeConfigure(configure string) string {
	configure = strings.TrimRight(configure, "\n")
	configure = strings.TrimRight(configure, " ")
	configure = strings.TrimRight(configure, "\\")
	if configure != "" {
		configure += " "
	}
	return configure
}

func normalizeAddModulePaths(path, rootDir string, dynamic bool) string {
	var result string
	if len(path) == 0 {
		return path
	}

	module_paths := strings.Split(path, ",")

	opt := "--add-module"
	if dynamic {
		opt = "--add-dynamic-module"
	}

	for _, module_path := range module_paths {
		if strings.HasPrefix(module_path, "/") {
			result += fmt.Sprintf("%s=%s \\\n", opt, module_path)
		} else {
			result += fmt.Sprintf("%s=%s/%s \\\n", opt, rootDir, module_path)
		}
	}

	return result
}

func configureNginx() error {
	args := []string{"sh", "./nginx-configure"}
	if VerboseEnabled {
		return runCommand(args)
	}

	f, err := os.Create("nginx-configure.log")
	if err != nil {
		return runCommand(args)
	}
	defer f.Close()

	cmd, err := makeCmd(args)
	if err != nil {
		return err
	}

	writer := bufio.NewWriter(f)
	cmd.Stdout = writer
	defer writer.Flush()

	return cmd.Run()
}
