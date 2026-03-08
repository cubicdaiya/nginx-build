package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	urlpkg "net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/cubicdaiya/nginx-build/builder"
	"github.com/cubicdaiya/nginx-build/command"
	"github.com/cubicdaiya/nginx-build/util"
)

const DefaultDownloadTimeout = time.Duration(900) * time.Second

func isGitURL(url string) bool {
	if url == "" {
		return false
	}

	// Check for explicit git transport forms first.
	if strings.HasPrefix(url, "git@") || strings.HasPrefix(url, "git://") {
		return true
	}

	parsedURL, err := urlpkg.Parse(url)
	if err != nil {
		return false
	}

	host := strings.ToLower(parsedURL.Hostname())
	urlPath := strings.Trim(parsedURL.EscapedPath(), "/")
	if urlPath == "" {
		return false
	}

	segments := strings.Split(urlPath, "/")
	lastSegment := segments[len(segments)-1]
	if strings.HasSuffix(lastSegment, ".git") {
		return true
	}

	if host == "codeload.github.com" || host == "api.github.com" {
		return false
	}

	if host == "github.com" {
		if len(segments) != 2 {
			return false
		}
		return true
	}

	if strings.HasSuffix(host, ".googlesource.com") || host == "googlesource.com" {
		for _, segment := range segments {
			if strings.HasPrefix(segment, "+archive") {
				return false
			}
		}
		return true
	}

	return false
}

func archiveEntryPath(name string) string {
	if name == "" {
		return ""
	}

	cleanPath := path.Clean(filepath.ToSlash(name))
	if cleanPath == "." {
		return ""
	}

	return filepath.FromSlash(cleanPath)
}

func extractArchive(path string) error {
	return command.Run([]string{"tar", "zxvf", path})
}

func archiveEntries(path string) ([]string, error) {
	cmd, err := command.Make([]string{"tar", "tzf", path})
	if err != nil {
		return nil, err
	}

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	entries := []string{}
	for _, line := range strings.Split(string(output), "\n") {
		entry := strings.TrimSpace(line)
		if entry == "" {
			continue
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func archiveRootDirectory(path string) (string, error) {
	entries, err := archiveEntries(path)
	if err != nil {
		return "", err
	}
	var root string
	for _, entry := range entries {
		entryPath := archiveEntryPath(entry)
		if entryPath == "" {
			continue
		}
		currentRoot := strings.SplitN(filepath.ToSlash(entryPath), "/", 2)[0]
		if currentRoot == "" || currentRoot == "." || currentRoot == ".." {
			continue
		}
		if root == "" {
			root = currentRoot
			continue
		}
		if root != currentRoot {
			return "", fmt.Errorf("archive %s must contain a single top-level directory", path)
		}
	}
	if root == "" {
		return "", fmt.Errorf("archive %s is empty", path)
	}
	return root, nil
}

func usesCustomSource(b *builder.Builder) bool {
	return b.CustomURL != ""
}

func usesGitDownload(b *builder.Builder) bool {
	return usesCustomSource(b) && isGitURL(b.DownloadURL())
}

func normalizeExtractedSourcePath(b *builder.Builder) error {
	if !usesCustomSource(b) || usesGitDownload(b) {
		return nil
	}

	rootDir, err := archiveRootDirectory(b.ArchivePath())
	if err != nil {
		return err
	}
	if rootDir == b.SourcePath() {
		return nil
	}
	if util.FileExists(b.SourcePath()) {
		return nil
	}

	info, err := os.Stat(rootDir)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("archive %s did not extract to a directory", b.ArchivePath())
	}

	return os.Rename(rootDir, b.SourcePath())
}

func download(b *builder.Builder) error {
	url := b.DownloadURL()

	if usesGitDownload(b) {
		// Clone from git
		log.Printf("Clone %s.....", b.SourcePath())
		if err := command.Run([]string{"git", "clone", url, b.SourcePath()}); err != nil {
			return fmt.Errorf("failed to clone from %s: %w", url, err)
		}

		// Checkout specific tag/branch if specified
		if b.CustomTag != "" {
			log.Printf("Checkout %s.....", b.CustomTag)
			originalDir, _ := os.Getwd()
			if err := os.Chdir(b.SourcePath()); err != nil {
				return fmt.Errorf("failed to change directory to %s: %w", b.SourcePath(), err)
			}
			if err := command.Run([]string{"git", "checkout", b.CustomTag}); err != nil {
				os.Chdir(originalDir)
				return fmt.Errorf("failed to checkout %s: %w", b.CustomTag, err)
			}
			if err := os.Chdir(originalDir); err != nil {
				return fmt.Errorf("failed to change back to original directory: %w", err)
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
		if usesGitDownload(b) {
			// Git clone handled in download()
			if err := download(b); err != nil {
				return fmt.Errorf("failed to download %s: %w", b.SourcePath(), err)
			}
		} else if !util.FileExists(b.ArchivePath()) {
			log.Printf("Download %s.....", b.SourcePath())

			if err := download(b); err != nil {
				return fmt.Errorf("failed to download %s: %w", b.SourcePath(), err)
			}
		}

		// Extract archive if it's not a git repository
		if !usesGitDownload(b) {
			log.Printf("Extract %s.....", b.ArchivePath())

			if err := extractArchive(b.ArchivePath()); err != nil {
				return fmt.Errorf("failed to extract %s: %w", b.ArchivePath(), err)
			}
			if err := normalizeExtractedSourcePath(b); err != nil {
				return fmt.Errorf("failed to prepare %s: %w", b.SourcePath(), err)
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
