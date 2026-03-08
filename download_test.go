package main

import (
	"archive/tar"
	"compress/gzip"
	"os"
	"path/filepath"
	"testing"
)

func TestIsGitURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want bool
	}{
		{
			name: "Git URL with .git suffix",
			url:  "https://github.com/user/repo.git",
			want: true,
		},
		{
			name: "Git URL with git:// protocol",
			url:  "git://github.com/user/repo",
			want: true,
		},
		{
			name: "GitHub URL without .git",
			url:  "https://github.com/user/repo",
			want: true,
		},
		{
			name: "GitHub releases URL",
			url:  "https://github.com/openssl/openssl/releases/download/openssl-3.5.1/openssl-3.5.1.tar.gz",
			want: false,
		},
		{
			name: "GitHub archive URL",
			url:  "https://github.com/user/repo/archive/v1.0.tar.gz",
			want: false,
		},
		{
			name: "GitHub tarball URL",
			url:  "https://github.com/user/repo/tarball/v1.0.0",
			want: false,
		},
		{
			name: "GitHub zipball URL",
			url:  "https://github.com/user/repo/zipball/v1.0.0",
			want: false,
		},
		{
			name: "GitHub codeload tarball URL",
			url:  "https://codeload.github.com/user/repo/tar.gz/refs/tags/v1.0.0",
			want: false,
		},
		{
			name: "GitHub non-repository page",
			url:  "https://github.com/user/repo/issues",
			want: false,
		},
		{
			name: "SSH GitHub URL",
			url:  "git@github.com:user/repo.git",
			want: true,
		},
		{
			name: "Google Source URL",
			url:  "https://boringssl.googlesource.com/boringssl",
			want: true,
		},
		{
			name: "Google Source archive URL",
			url:  "https://boringssl.googlesource.com/boringssl/+archive/refs/heads/main.tar.gz",
			want: false,
		},
		{
			name: "Regular tar.gz URL",
			url:  "https://example.com/library.tar.gz",
			want: false,
		},
		{
			name: "Regular tar.bz2 URL",
			url:  "https://example.com/library.tar.bz2",
			want: false,
		},
		{
			name: "Empty URL",
			url:  "",
			want: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := isGitURL(test.url); got != test.want {
				t.Errorf("isGitURL(%q) = %v, want %v", test.url, got, test.want)
			}
		})
	}
}

func TestArchiveRootDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	archivePath := filepath.Join(tmpDir, "source.tar.gz")

	f, err := os.Create(archivePath)
	if err != nil {
		t.Fatalf("failed to create archive: %v", err)
	}

	gz := gzip.NewWriter(f)
	tw := tar.NewWriter(gz)
	headers := []tar.Header{
		{Name: "custom-nginx/", Mode: 0755, Typeflag: tar.TypeDir},
		{Name: "custom-nginx/configure", Mode: 0644, Size: int64(len("content"))},
	}
	for _, hdr := range headers {
		h := hdr
		if err := tw.WriteHeader(&h); err != nil {
			t.Fatalf("failed to write header: %v", err)
		}
		if hdr.Size > 0 {
			if _, err := tw.Write([]byte("content")); err != nil {
				t.Fatalf("failed to write content: %v", err)
			}
		}
	}
	if err := tw.Close(); err != nil {
		t.Fatalf("failed to close tar writer: %v", err)
	}
	if err := gz.Close(); err != nil {
		t.Fatalf("failed to close gzip writer: %v", err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("failed to close archive: %v", err)
	}

	rootDir, err := archiveRootDirectory(archivePath)
	if err != nil {
		t.Fatalf("archiveRootDirectory() returned error: %v", err)
	}
	if rootDir != "custom-nginx" {
		t.Fatalf("archiveRootDirectory() = %q, want %q", rootDir, "custom-nginx")
	}
}

func TestExtractArchiveDoesNotWriteOutsideDestination(t *testing.T) {
	tmpDir := t.TempDir()
	archivePath := filepath.Join(tmpDir, "traversal.tar.gz")
	extractDir := filepath.Join(tmpDir, "extract")
	if err := os.Mkdir(extractDir, 0755); err != nil {
		t.Fatalf("failed to create extract dir: %v", err)
	}

	f, err := os.Create(archivePath)
	if err != nil {
		t.Fatalf("failed to create archive: %v", err)
	}

	gz := gzip.NewWriter(f)
	tw := tar.NewWriter(gz)
	content := "content"
	if err := tw.WriteHeader(&tar.Header{
		Name:     "safe/",
		Mode:     0755,
		Typeflag: tar.TypeDir,
	}); err != nil {
		t.Fatalf("failed to write dir header: %v", err)
	}
	if err := tw.WriteHeader(&tar.Header{
		Name:     "../evil.txt",
		Mode:     0644,
		Size:     int64(len(content)),
		Typeflag: tar.TypeReg,
	}); err != nil {
		t.Fatalf("failed to write file header: %v", err)
	}
	if _, err := tw.Write([]byte(content)); err != nil {
		t.Fatalf("failed to write file content: %v", err)
	}
	if err := tw.Close(); err != nil {
		t.Fatalf("failed to close tar writer: %v", err)
	}
	if err := gz.Close(); err != nil {
		t.Fatalf("failed to close gzip writer: %v", err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("failed to close archive: %v", err)
	}

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}
	if err := os.Chdir(extractDir); err != nil {
		t.Fatalf("failed to change dir: %v", err)
	}
	defer os.Chdir(wd)

	err = extractArchive(archivePath)
	if _, statErr := os.Stat(filepath.Join(tmpDir, "evil.txt")); !os.IsNotExist(statErr) {
		t.Fatalf("unexpected file extracted outside target directory: %v", statErr)
	}
	if err == nil {
		if _, statErr := os.Stat(filepath.Join(extractDir, "evil.txt")); statErr != nil {
			t.Fatalf("expected tar to keep extracted file within destination directory when no error is returned: %v", statErr)
		}
	}
}
