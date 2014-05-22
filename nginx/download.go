package nginx

import (
	"fmt"
	"os/exec"
)

func DownloadLink(version string) string {
	return fmt.Sprintf("%s/%s", NGINX_DOWNLOAD_URL_PREFIX, ArchivePath(version))
}

func Download(link string) error {
	cmd := exec.Command("wget", link)
	checkVerboseenabled(cmd)
	return cmd.Run()
}

func DownloadModule3rd(module3rd Module3rd) error {
	form := module3rd.form
	url := module3rd.url
	var cmd *exec.Cmd
	switch form {
	case "github":
		cmd = exec.Command("git", "clone", url)
	}
	checkVerboseenabled(cmd)
	return cmd.Run()
}
