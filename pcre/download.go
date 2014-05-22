package pcre

import (
	"../common"
	"fmt"
	"os/exec"
)

func DownloadLink(version string) string {
	return fmt.Sprintf("%s/%s", DOWNLOAD_URL_PREFIX, ArchivePath(version))
}

func Download(link string) error {
	cmd := exec.Command("wget", link)
	common.CheckVerboseEnabled(cmd)
	return cmd.Run()
}
