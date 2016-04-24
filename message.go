package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
)

func printNginxBuildVersion() {
	fmt.Printf(`nginx-build %s
Compiler: %s %s
Copyright (C) 2014-2016 Tatsuhiko Kubo <cubicdaiya@gmail.com>
`,
		NGINX_BUILD_VERSION,
		runtime.Compiler,
		runtime.Version())

}

func printConfigureOptions() error {
	cmd := exec.Command("objs/nginx", "-V")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func printFirstMsg() {
	fmt.Printf(`nginx-build: %s
Compiler: %s %s
`,
		NGINX_BUILD_VERSION,
		runtime.Compiler,
		runtime.Version())
}

func printLastMsg(workDir, srcDir string, openResty, configureOnly bool) {
	log.Println("Complete building nginx!")

	if !openResty {
		if !configureOnly {
			fmt.Println()
			err := printConfigureOptions()
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}
	fmt.Println()

	lastMsgFormat := `Enter the following command for install nginx.

   $ cd %s/%s%s
   $ sudo make install
`
	if configureOnly {
		log.Printf(lastMsgFormat, workDir, srcDir, "\n   $ make")
	} else {
		log.Printf(lastMsgFormat, workDir, srcDir, "")
	}
}
