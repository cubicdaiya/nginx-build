package main

import (
	"fmt"
	"os/exec"
)

func (builder *Builder) DownloadLink() string {
	return fmt.Sprintf("%s/%s", builder.DownLoadPrefix, builder.ArchivePath())
}

func (builder *Builder) Download() error {
	link := builder.DownloadLink()
	cmd := exec.Command("wget", link)
	CheckVerboseEnabled(cmd)
	return cmd.Run()
}

func (builder *Builder) DownloadModule3rd(module3rd Module3rd) error {
	form := module3rd.Form
	url := module3rd.Url
	var cmd *exec.Cmd
	switch form {
	case "github":
		cmd = exec.Command("git", "clone", url)
		CheckVerboseEnabled(cmd)
		return cmd.Run()
	}
	return nil
}
