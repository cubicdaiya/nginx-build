package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/cubicdaiya/nginx-build/builder"
	"github.com/cubicdaiya/nginx-build/command"
	"github.com/cubicdaiya/nginx-build/util"
)

func extractArchive(path string) error {
	return command.Run([]string{"tar", "zxvf", path})
}

func download(b *builder.Builder) error {
	c := &http.Client{
		Timeout: b.DownloadTimeout,
	}
	res, err := c.Get(b.DownloadURL())
	if err != nil {
		return err
	}
	defer res.Body.Close()

	tmpFileName := b.ArchivePath() + ".download"
	f, err := os.Create(tmpFileName)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := io.Copy(f, res.Body); err != nil && err != io.EOF {
		return err
	}

	if err := os.Rename(tmpFileName, b.ArchivePath()); err != nil {
		return err
	}

	return nil
}

func downloadAndExtract(b *builder.Builder) error {
	if !util.FileExists(b.SourcePath()) {
		if !util.FileExists(b.ArchivePath()) {

			log.Printf("Download %s.....", b.SourcePath())

			if err := download(b); err != nil {
				return fmt.Errorf("Failed to download %s. %s", b.SourcePath(), err.Error())
			}
		}

		log.Printf("Extract %s.....", b.ArchivePath())

		if err := extractArchive(b.ArchivePath()); err != nil {
			return fmt.Errorf("Failed to extract %s. %s", b.ArchivePath(), err.Error())
		}
	} else {
		log.Printf("%s already exists.", b.SourcePath())
	}
	return nil
}

func downloadAndExtractParallel(b *builder.Builder) {
	if err := downloadAndExtract(b); err != nil {
		util.PrintFatalMsg(err, b.LogPath())
	}
}
