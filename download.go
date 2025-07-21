package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/cubicdaiya/nginx-build/builder"
	"github.com/cubicdaiya/nginx-build/command"
	"github.com/cubicdaiya/nginx-build/util"
)

const DefaultDownloadTimeout = time.Duration(900) * time.Second

func isGitURL(url string) bool {
	// Check for explicit git URLs
	if strings.HasSuffix(url, ".git") || strings.Contains(url, "git://") {
		return true
	}

	// Check for git hosting services, but exclude release/download URLs
	if strings.Contains(url, "github.com") && !strings.Contains(url, "/releases/download/") && !strings.Contains(url, "/archive/") {
		return true
	}

	if strings.Contains(url, "googlesource.com") {
		return true
	}

	return false
}

func extractArchive(path string) error {
	return command.Run([]string{"tar", "zxvf", path})
}

func download(b *builder.Builder) error {
	url := b.DownloadURL()

	// Only check for git URLs if this is a custom SSL component
	if b.Component == builder.ComponentCustomSSL && url != "" && isGitURL(url) {
		// Clone from git
		log.Printf("Clone %s.....", b.SourcePath())
		if err := command.Run([]string{"git", "clone", url, b.SourcePath()}); err != nil {
			return fmt.Errorf("Failed to clone from %s. %s", url, err.Error())
		}

		// Checkout specific tag/branch if specified
		if b.CustomTag != "" {
			log.Printf("Checkout %s.....", b.CustomTag)
			originalDir, _ := os.Getwd()
			if err := os.Chdir(b.SourcePath()); err != nil {
				return fmt.Errorf("Failed to change directory to %s. %s", b.SourcePath(), err.Error())
			}
			if err := command.Run([]string{"git", "checkout", b.CustomTag}); err != nil {
				os.Chdir(originalDir)
				return fmt.Errorf("Failed to checkout %s. %s", b.CustomTag, err.Error())
			}
			if err := os.Chdir(originalDir); err != nil {
				return fmt.Errorf("Failed to change back to original directory. %s", err.Error())
			}
		}

		return nil
	}

	c := &http.Client{
		Timeout: DefaultDownloadTimeout,
	}
	res, err := c.Get(b.DownloadURL())
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download %s. %s", b.DownloadURL(), res.Status)
	}

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
		// For custom SSL or other components
		if b.Component == builder.ComponentCustomSSL && isGitURL(b.DownloadURL()) {
			// Git clone handled in download()
			if err := download(b); err != nil {
				return fmt.Errorf("Failed to download %s. %s", b.SourcePath(), err.Error())
			}
		} else if !util.FileExists(b.ArchivePath()) {
			log.Printf("Download %s.....", b.SourcePath())

			if err := download(b); err != nil {
				return fmt.Errorf("Failed to download %s. %s", b.SourcePath(), err.Error())
			}
		}

		// Extract archive if it's not a git repository
		if !(b.Component == builder.ComponentCustomSSL && isGitURL(b.DownloadURL())) {
			log.Printf("Extract %s.....", b.ArchivePath())

			if err := extractArchive(b.ArchivePath()); err != nil {
				return fmt.Errorf("Failed to extract %s. %s", b.ArchivePath(), err.Error())
			}
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
