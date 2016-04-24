package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/cubicdaiya/nginx-build/util"
)

func clearWorkDir(workDir string) error {
	err := os.RemoveAll(workDir)
	if err != nil {
		// workaround for the restriction of os.RemoveAll()
		// os.RemoveAll() calls fd.Readdirnames(100).
		// So os.RemoveAll() does not always remove all entries.
		// Some 3rd-party module (e.g. lua-nginx-module) tumbles this restriction.
		if util.FileExists(workDir) {
			err = os.RemoveAll(workDir)
		}
	}
	return err
}

func fileGetContents(path string) (string, error) {
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
