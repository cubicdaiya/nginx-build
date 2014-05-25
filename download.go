package main

import (
	"fmt"
	"log"
	"os/exec"
)

func download(url string) error {
	return runCommand(exec.Command("wget", url))
}

func downloadModule3rd(module3rd Module3rd) error {
	form := module3rd.Form
	url := module3rd.Url
	switch form {
	case "git":
		return runCommand(exec.Command("git", "clone", url))
	}
	return nil
}

func downloadAndExtract(builder *Builder) error {
	if !fileExists(builder.sourcePath()) {
		name := builder.name()
		if !fileExists(builder.archivePath()) {
			log.Printf("Download %s.....", name)
			url := builder.downloadURL()
			err := download(url)
			if err != nil {
				return fmt.Errorf("Failed to download %s", name)
			}
		}
		log.Printf("Extract %s.....", name)
		err := extractArchive(builder.archivePath())
		if err != nil {
			return fmt.Errorf("Failed to extract %s", name)
		}
	} else {
		log.Printf("%s already exists.", builder.sourcePath())
	}
	return nil
}
