package nginx

import (
	"../common"
	"os/exec"
	"runtime"
	"strconv"
)

func Configure() error {
	cmd := exec.Command("sh", "./nginx-configure")
	common.CheckVerboseEnabled(cmd)
	return cmd.Run()
}

func Make(conf string) error {
	numCPU := runtime.NumCPU()
	cmd := exec.Command("make", "-j", strconv.Itoa(numCPU))
	common.CheckVerboseEnabled(cmd)
	return cmd.Run()
}

func ExtractArchive(path string) error {
	cmd := exec.Command("tar", "zxvf", path)
	common.CheckVerboseEnabled(cmd)
	return cmd.Run()
}
