package util

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/cubicdaiya/nginx-build/command"
)

var (
	mutex   sync.Mutex
	patched bool
)

func init() {
	patched = false
}

func patch(path, option string, reverse bool) error {

	if reverse {
		mutex.Lock()
		if patched {
			mutex.Unlock()
			return nil
		}
		patched = true
		mutex.Unlock()
	}

	args := []string{"sh", "-c"}
	body := ""
	if reverse {
		body = fmt.Sprintf("patch %s -R < %s", option, path)
	} else {
		body = fmt.Sprintf("patch %s < %s", option, path)
	}
	args = append(args, body)

	cmd, err := command.Make(args)
	if err != nil {
		return err
	}

	// As the output of patch is interactive,
	// the result is always output.
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func Patch(pathInput, option, rootDir string, reverse bool) error {
	if pathInput == "" {
		return nil
	}

	var individualPaths []string
	if strings.Contains(pathInput, ",") {
		individualPaths = strings.Split(pathInput, ",")
	} else {
		individualPaths = append(individualPaths, pathInput)
	}

	var expandedPaths []string
	for _, p := range individualPaths {
		var currentPath string
		if filepath.IsAbs(p) {
			currentPath = p
		} else {
			currentPath = filepath.Join(rootDir, p)
		}

		isDir, err := IsDirectory(currentPath)
		if err != nil {
			return fmt.Errorf("error checking if patch path %s is directory: %w", currentPath, err)
		}
		if isDir {
			pathsInDir, err := ListDirectory(currentPath)
			if err != nil {
				return fmt.Errorf("error listing directory for patches %s: %w", currentPath, err)
			}
			if pathsInDir != nil {
				expandedPaths = append(expandedPaths, pathsInDir...)
			}
		} else {
			expandedPaths = append(expandedPaths, currentPath)
		}
	}

	for _, p := range expandedPaths {
		if !FileExists(p) {
			return fmt.Errorf("patch file %s not found", p)
		}

		logMsg := "Applying"
		if reverse {
			logMsg = "Reverting"
		}
		log.Printf("%s patch: %s (options: %s)", logMsg, p, option)

		if err := patch(p, option, reverse); err != nil {
			return fmt.Errorf("failed to %s patch %s (options: %s): %w", strings.ToLower(logMsg), p, option, err)
		}
	}
	return nil
}
