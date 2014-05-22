package nginx

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

func DownloadModule3rd(module3rd Module3rd) error {
	form := module3rd.Form
	url := module3rd.Url
	var cmd *exec.Cmd
	switch form {
	case "github":
		cmd = exec.Command("git", "clone", url)
		common.CheckVerboseEnabled(cmd)
		return cmd.Run()
	}
	return nil
}
