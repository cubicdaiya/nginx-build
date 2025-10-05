package configure

import (
	"strings"
)

// Normalize prepares the base configure script and captures trailing comment lines.
func Normalize(configure string) (string, string) {
	lines := strings.Split(configure, "\n")

	// Drop trailing blank lines.
	for len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
		lines = lines[:len(lines)-1]
	}

	// Capture trailing comment lines (prefixed with # after trimming spaces).
	var trailingComments []string
	for len(lines) > 0 {
		trimmed := strings.TrimSpace(lines[len(lines)-1])
		if strings.HasPrefix(trimmed, "#") {
			trailingComments = append([]string{lines[len(lines)-1]}, trailingComments...)
			lines = lines[:len(lines)-1]
			continue
		}
		break
	}

	base := strings.Join(lines, "\n")
	base = strings.TrimRight(base, "\n")
	base = strings.TrimRight(base, " ")
	base = strings.TrimRight(base, "\\")
	if base != "" {
		base += " "
	}

	return base, strings.Join(trailingComments, "\n")
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
			result += opt + "=" + modulePath + " \\\n"
		} else {
			result += opt + "=" + rootDir + "/" + modulePath + " \\\n"
		}
	}

	return result
}
