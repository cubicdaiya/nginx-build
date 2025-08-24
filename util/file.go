package util

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
}

func IsDirectory(path string) (bool, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return stat.IsDir(), nil
}

func ListDirectory(root string) ([]string, error) {
	var files []string
	var walkErr error

	err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			if errors.Is(err, fs.ErrPermission) && d != nil && d.IsDir() {
				return fs.SkipDir
			}
			walkErr = errors.Join(walkErr, fmt.Errorf("%s: %w", p, err))
			return nil
		}
		if d != nil && !d.IsDir() {
			files = append(files, p)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, walkErr
}

func SaveCurrentDir() (string, error) {
	prevDir, err := filepath.Abs(".")
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path for current directory: %w", err)
	}
	return prevDir, nil
}

func FileGetContents(path string) (string, error) {
	conf := ""
	if len(path) > 0 {
		confb, err := os.ReadFile(path)
		if err != nil {
			return "", fmt.Errorf("confPath(%s) does not exist.", path)
		}
		conf = string(confb)
	}
	return conf, nil
}

func ClearWorkDir(workDir string) error {
	err := os.RemoveAll(workDir)
	if err != nil {
		// workaround for the restriction of os.RemoveAll()
		// os.RemoveAll() calls fd.Readdirnames(100).
		// So os.RemoveAll() does not always remove all entries.
		// Some 3rd-party module (e.g. lua-nginx-module) tumbles this restriction.
		if FileExists(workDir) {
			err = os.RemoveAll(workDir)
		}
	}
	return err
}
