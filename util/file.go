package util

import (
	"fmt"
	"io/ioutil"
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

func ListDirectory(path string) ([]string, error) {
	var fileNames []string

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			fileNames = append(fileNames, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return fileNames, nil
}

func SaveCurrentDir() string {
	prevDir, _ := filepath.Abs(".")
	return prevDir
}

func FileGetContents(path string) (string, error) {
	conf := ""
	if len(path) > 0 {
		confb, err := ioutil.ReadFile(path)
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
