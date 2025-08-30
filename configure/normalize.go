package configure

import (
	"fmt"
	"strings"
)

func Normalize(configure string) string {
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

	modulePaths := strings.Split(path, ",")

	opt := "--add-module"
	if dynamic {
		opt = "--add-dynamic-module"
	}

	for _, modulePath := range modulePaths {
		if strings.HasPrefix(modulePath, "/") {
			result += fmt.Sprintf("%s=%s \\\n", opt, modulePath)
		} else {
			result += fmt.Sprintf("%s=%s/%s \\\n", opt, rootDir, modulePath)
		}
	}

	return result
}
