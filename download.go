package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/cubicdaiya/nginx-build/builder"
	"github.com/cubicdaiya/nginx-build/command"
	"github.com/cubicdaiya/nginx-build/util"
)

const DefaultDownloadTimeout = time.Duration(900) * time.Second

func extractArchive(path string) error {
	return command.Run([]string{"tar", "zxvf", path})
}

func download(b *builder.Builder) error {
	c := &http.Client{
		Timeout: DefaultDownloadTimeout,
	}

	res, err := c.Get(b.DownloadURL())
	if err != nil {
		return err
	}
	defer res.Body.Close()

	f, err := os.Create(b.ArchivePath())
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, res.Body)
	if err != nil && err != io.EOF {
		return err
	}

	return nil
}

func downloadAndExtract(b *builder.Builder) error {
	if !util.FileExists(b.SourcePath()) {
		if !util.FileExists(b.ArchivePath()) {

			log.Printf("Download %s.....", b.SourcePath())

			err := download(b)
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

func downloadAndExtractParallel(b *builder.Builder) {
	err := downloadAndExtract(b)
	if err != nil {
		util.PrintFatalMsg(err, b.LogPath())
	}
}
