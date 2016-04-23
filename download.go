package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/cubicdaiya/nginx-build/builder"
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
		fallthrough
	case "hg":
		args := []string{form, "clone", url}
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

func downloadAndExtract(b *builder.Builder) error {
	if !fileExists(b.SourcePath()) {
		if !fileExists(b.ArchivePath()) {

			log.Printf("Download %s.....", b.SourcePath())

			err := download(b.DownloadURL(), b.LogPath())
			if err != nil {
				return fmt.Errorf("Failed to download %s. %s", b.SourcePath(), err.Error())
			}
		}

		log.Printf("Extract %s.....", b.ArchivePath())

		err := extractArchive(b.ArchivePath())
		if err != nil {
			return fmt.Errorf("Failed to extract %s. %s", b.ArchivePath(), err.Error())
		}
	} else {
		log.Printf("%s already exists.", b.SourcePath())
	}
	return nil
}

func downloadAndExtractParallel(b *builder.Builder, wg *sync.WaitGroup) {
	err := downloadAndExtract(b)
	if err != nil {
		printFatalMsg(err, b.LogPath())
	}
	wg.Done()
}

func downloadAndExtractModule3rdParallel(m Module3rd, wg *sync.WaitGroup) {
	if fileExists(m.Name) {
		log.Printf("%s already exists.", m.Name)
		wg.Done()
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

	wg.Done()
}
