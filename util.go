package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
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

func printFirstMsg() {
	fmt.Printf(`nginx-build: %s
Compiler: %s %s
`,
		NGINX_BUILD_VERSION,
		runtime.Compiler,
		runtime.Version())
}

func printLastMsg(workDir, srcDir string, openResty bool) {
	log.Println("Complete building nginx!")

	fmt.Println()
	if !openResty {
		err := printConfigureOptions()
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	fmt.Println()

	lastMsgFormat := `Enter the following command for install nginx.

   $ cd %s/%s
   $ sudo make install
`
	log.Printf(lastMsgFormat, workDir, srcDir)
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

func configureNginx() error {
	return runCommand(exec.Command("sh", "./nginx-configure"))
}

func buildNginx(jobs int) error {
	return runCommand(exec.Command("make", "-j", strconv.Itoa(jobs)))
}

func extractArchive(path string) error {
	return runCommand(exec.Command("tar", "zxvf", path))
}

func switchRev(rev string) error {
	return runCommand(exec.Command("git", "checkout", rev))
}

func printConfigureOptions() error {
	cmd := exec.Command("objs/nginx", "-V")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func provideShell(sh string) error {
	if len(sh) == 0 {
		return nil
	}
	args := strings.Split(strings.Trim(sh, " "), " ")
	var err error
	if len(args) == 1 {
		err = runCommand(exec.Command(args[0]))
	} else {
		err = runCommand(exec.Command(args[0], args[1:]...))
	}
	return err
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
