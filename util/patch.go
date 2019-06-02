package util

import (
	"fmt"
	"log"
	"os"
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

func Patch(path, option, root string, reverse bool) {
	if path == "" {
		return
	}

	var pathes []string
	if strings.Contains(path, ",") {
		pathes = strings.Split(path, ",")
	} else {
		pathes = append(pathes, path)
	}

	// replace directories with all files they contain (recursively)
	var expanded_paths []string
	for _, path := range pathes {
		if !strings.HasPrefix(path, "/") {
			path = fmt.Sprintf("%s/%s", root, path)
		}

		isDir, err := IsDirectory(path)
		if err != nil {
			log.Fatal(err)
		}
		if isDir {
			paths, err := ListDirectory(path)
			if err != nil {
				log.Fatal(err)
			}
			if paths != nil {
				expanded_paths = append(expanded_paths, paths...)
			}
		} else {
			expanded_paths = append(expanded_paths, path)
		}
	}

	pathes = expanded_paths

	for _, path := range pathes {
		if FileExists(path) {
			if reverse {
				log.Printf("Reverting patch: %s", path)
			} else {
				log.Printf("Applying patch: %s %s", option, path)
			}
			if err := patch(path, option, reverse); err != nil {
				log.Fatalf("Failed to apply patch: %s %s", option, path)
			}
		} else {
			log.Fatalf("Patch pathname: %s is not found", path)
		}
	}
}
