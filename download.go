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
	return fmt.Errorf("form=%s is not supported", form)
}

func downloadAndExtract(builder *Builder) error {
	if !fileExists(builder.sourcePath()) {
		if !fileExists(builder.archivePath()) {
			log.Printf("Download %s.....", builder.sourcePath())
			url := builder.downloadURL()
			err := download(url)
			if err != nil {
				return fmt.Errorf("Failed to download %s", builder.sourcePath())
			}
		}
		log.Printf("Extract %s.....", builder.archivePath())
		err := extractArchive(builder.archivePath())
		if err != nil {
			return fmt.Errorf("Failed to extract %s", builder.archivePath())
		}
	} else {
		log.Printf("%s already exists.", builder.sourcePath())
	}
	return nil
}

func downloadAndExtractParallel(builder *Builder, done chan bool) {
	err := downloadAndExtract(builder)
	if err != nil {
		log.Fatal(err.Error())
	}
	done <- true
}

func downloadAndExtractModule3rdParallel(m Module3rd, done chan bool) {
	if fileExists(m.Name) {
		log.Printf("%s already exists.", m.Name)
		done <- true
		return
	}
	if len(m.Rev) > 0 {
		log.Printf("Download %s-%s.....", m.Name, m.Rev)
	} else {
		log.Printf("Download %s.....", m.Name)
	}
	err := downloadModule3rd(m)
	if err != nil {
		log.Println(err.Error())
		log.Fatalf("Failed to download %s", m.Name)
	}
	done <- true
}
