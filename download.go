package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"sync"

	"github.com/cubicdaiya/nginx-build/builder"
	"github.com/cubicdaiya/nginx-build/command"
	"github.com/cubicdaiya/nginx-build/util"
)

func extractArchive(path string) error {
	return command.Run([]string{"tar", "zxvf", path})
}

func downloadBuiltin(url string, logName string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept-Encoding", "gzip")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	filename := path.Base(url)
	if err := ioutil.WriteFile(filename, body, 0644); err != nil {
		return err
	}

	return nil
}

func downloadByWget(url string, logName string) error {
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

func downloadAndExtract(b *builder.Builder) error {
	if !util.FileExists(b.SourcePath()) {
		if !util.FileExists(b.ArchivePath()) {

			log.Printf("Download %s.....", b.SourcePath())

			var err error
			if command.VerboseEnabled {
				err = downloadByWget(b.DownloadURL(), b.LogPath())
			} else {
				err = downloadBuiltin(b.DownloadURL(), b.LogPath())
			}
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
