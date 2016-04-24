package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/cubicdaiya/nginx-build/builder"
	"github.com/cubicdaiya/nginx-build/command"
	"github.com/cubicdaiya/nginx-build/module3rd"
	"github.com/cubicdaiya/nginx-build/util"
)

func extractArchive(path string) error {
	return command.Run([]string{"tar", "zxvf", path})
}

func download(url string, logName string) error {
	args := []string{"wget", url}
	if command.VerboseEnabled {
		return command.Run(args)
	}

	f, err := os.Create(logName)
	if err != nil {
		return command.Run(args)
	}
	defer f.Close()

	cmd, err := command.Make(args)
	if err != nil {
		return err
	}

	writer := bufio.NewWriter(f)
	defer writer.Flush()

	cmd.Stderr = writer

	return cmd.Run()
}

func downloadModule3rd(m module3rd.Module3rd, logName string) error {
	form := m.Form
	url := m.Url

	switch form {
	case "git":
		fallthrough
	case "hg":
		args := []string{form, "clone", url}
		if command.VerboseEnabled {
			return command.Run(args)
		}

		f, err := os.Create(logName)
		if err != nil {
			return command.Run(args)
		}
		defer f.Close()

		cmd, err := command.Make(args)
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
	if !util.FileExists(b.SourcePath()) {
		if !util.FileExists(b.ArchivePath()) {

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
		util.PrintFatalMsg(err, b.LogPath())
	}
	wg.Done()
}

func downloadAndExtractModule3rdParallel(m module3rd.Module3rd, wg *sync.WaitGroup) {
	if util.FileExists(m.Name) {
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
			util.PrintFatalMsg(err, logName)
		}
	} else if !util.FileExists(m.Url) {
		log.Fatalf("no such directory:%s", m.Url)
	}

	wg.Done()
}
