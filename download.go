package main

import (
	"fmt"
	"log"
	"bufio"
	"os"
)

func extractArchive(path string) error {
	return runCommand([]string{"tar", "zxvf", path})
}

func download(url string, logName string) error {
	args := []string{"wget", url}
	if VerboseEnabled {
		return runCommand(args)
	}

	f, err := os.Create(logName)
	if err != nil {
		return runCommand(args)
	}
	defer f.Close()

	cmd, err := makeCmd(args)
	if err != nil {
		return err
	}

	writer := bufio.NewWriter(f)
	defer writer.Flush()

	cmd.Stderr = writer

	return cmd.Run()
}

func downloadModule3rd(module3rd Module3rd, logName string) error {
	form := module3rd.Form
	url := module3rd.Url

	switch form {
	case "git":
		args := []string{"git", "clone", url}
		if VerboseEnabled {
			return runCommand(args)
		}

		f, err := os.Create(logName)
		if err != nil {
			return runCommand(args)
		}
		defer f.Close()

		cmd, err := makeCmd(args)
		if err != nil {
			return err
		}

		writer := bufio.NewWriter(f)
		defer writer.Flush()

		cmd.Stderr = writer

		return cmd.Run()
	case "local":
		return nil
	}

	return fmt.Errorf("form=%s is not supported", form)
}

func downloadAndExtract(builder *Builder) error {
	if !fileExists(builder.sourcePath()) {
		if !fileExists(builder.archivePath()) {

			log.Printf("Download %s.....", builder.sourcePath())

			err := download(builder.downloadURL(), builder.logPath())
			if err != nil {
				return fmt.Errorf("Failed to download %s. %s", builder.sourcePath(), err.Error())
			}
		}

		log.Printf("Extract %s.....", builder.archivePath())

		err := extractArchive(builder.archivePath())
		if err != nil {
			return fmt.Errorf("Failed to extract %s. %s", builder.archivePath(), err.Error())
		}
	} else {
		log.Printf("%s already exists.", builder.sourcePath())
	}
	return nil
}

func downloadAndExtractParallel(builder *Builder, done chan bool) {
	err := downloadAndExtract(builder)
	if err != nil {
		printFatalMsg(err, builder.logPath())
	}
	done <- true
}

func downloadAndExtractModule3rdParallel(m Module3rd, done chan bool) {
	if fileExists(m.Name) {
		log.Printf("%s already exists.", m.Name)
		done <- true
		return
	}

	if m.Form != "local" {
		if len(m.Rev) > 0 {
			log.Printf("Download %s-%s.....", m.Name, m.Rev)
		} else {
			log.Printf("Download %s.....", m.Name)
		}

		logName := fmt.Sprintf("%s.log", m.Name)

		err := downloadModule3rd(m, logName)
		if err != nil {
			printFatalMsg(err, logName)
		}
	} else if !fileExists(m.Url) {
		log.Fatalf("no such directory:%s", m.Url)
	}

	done <- true
}
