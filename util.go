package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func runCommand(cmd *exec.Cmd) error {
	checkVerboseEnabled(cmd)
	return cmd.Run()
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
}

func saveCurrentDir() string {
	prevDir, _ := filepath.Abs(".")
	return prevDir
}

func restoreCurrentDir(prevDir string) {
	os.Chdir(prevDir)
}

func printLastMsg(workDir, srcDir string) {
	log.Println("Complete building nginx!")

	lastMsgFormat := `Enter the following command for install nginx.

   $ cd %s/%s
   $ sudo make install
`
	log.Println(fmt.Sprintf(lastMsgFormat, workDir, srcDir))
}

func versionCheck(version string) {
	if len(version) == 0 {
		log.Println("[warn]nginx version is not set.")
		log.Printf("[warn]nginx-build use %s.\n", NGINX_VERSION)
	}
}

func fileGetContents(path string) (string, error) {
	conf := ""
	if len(path) == 0 {
		log.Println("[warn]configure option is empty.")
	} else {
		confb, err := ioutil.ReadFile(path)
		if err != nil {
			return "", fmt.Errorf("confPath(%s) does not exist.", path)
		}
		conf = string(confb)
	}
	return conf, nil
}
