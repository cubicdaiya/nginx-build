package pcre

import (
	"../common"
	"os/exec"
)

func ExtractArchive(path string) error {
	cmd := exec.Command("tar", "zxvf", path)
	common.CheckVerboseEnabled(cmd)
	return cmd.Run()
}
