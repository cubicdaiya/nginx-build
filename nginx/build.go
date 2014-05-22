package nginx

import (
	"../common"
	"os/exec"
	"strconv"
)

func Configure() error {
	cmd := exec.Command("sh", "./nginx-configure")
	common.CheckVerboseEnabled(cmd)
	return cmd.Run()
}

func Make(conf string, jobs int) error {
	cmd := exec.Command("make", "-j", strconv.Itoa(jobs))
	common.CheckVerboseEnabled(cmd)
	return cmd.Run()
}

func ExtractArchive(path string) error {
	cmd := exec.Command("tar", "zxvf", path)
	common.CheckVerboseEnabled(cmd)
	return cmd.Run()
}

func SwitchRev(rev string) error {
	cmd := exec.Command("git", "co", rev)
	common.CheckVerboseEnabled(cmd)
	return cmd.Run()
}
