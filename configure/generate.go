package configure

import (
	"fmt"
	"strings"

	"github.com/cubicdaiya/nginx-build/builder"
	"github.com/cubicdaiya/nginx-build/module3rd"
)

func Generate(configure string, modules3rd []module3rd.Module3rd, dependencies []builder.StaticLibrary, options Options, rootDir string, openResty bool, jobs int) string {
	openSSLStatic := false
	if len(configure) == 0 {
		configure = `#!/bin/sh

./configure \
`
	}

	if openResty {
		configure += fmt.Sprintf("-j%d \\\n", jobs)
	}

	for _, d := range dependencies {
		configure += fmt.Sprintf("%s=../%s-%s \\\n", d.Option, d.Name, d.Version)
		if d.Name == "openssl" || d.Name == "libressl" {
			openSSLStatic = true
		}
	}

	if openSSLStatic && !strings.Contains(configure, "--with-http_ssl_module") {
		configure += "--with-http_ssl_module \\\n"
	}

	configureModules3rd := generateForModule3rd(modules3rd)
	configure += configureModules3rd

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

func generateForModule3rd(modules3rd []module3rd.Module3rd) string {
	var result strings.Builder
	for _, m := range modules3rd {
		opt := "--add-module"
		if m.Dynamic {
			opt = "--add-dynamic-module"
		}
		if m.Form == "local" {
			result.WriteString(fmt.Sprintf("%s=%s \\\n", opt, m.URL))
		} else {
			result.WriteString(fmt.Sprintf("%s=../%s \\\n", opt, m.Name))
		}
	}
	return result.String()
}
