package main

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func setupTestArchive(t *testing.T, src string) {

	dst, err := os.Create("tests/test1.tar.gz")
	if err != nil {
		panic(err)
	}
	defer dst.Close()

	gw := gzip.NewWriter(dst)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	err = filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		h, err := tar.FileInfoHeader(info, path)
		if err != nil {
			return err
		}

		// tar.FileInfoHeader() assigns info.Name() to h.Name
		// overwrite with the necessary path
		h.Name = path

		if err := tw.WriteHeader(h); err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		if _, err := io.Copy(tw, f); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		panic(err)
	}
}

func TestUntargz(t *testing.T) {
	setupTestArchive(t, "tests/test1")

	tests := []struct {
		dst string
		src string
	}{
		{
			dst: "tests/test1_result",
			src: "tests/test1.tar.gz",
		},
	}

	for _, test := range tests {
		if err := untargz(test.dst, test.src); err != nil {
			t.Fatalf("Failed to un-tar %v to %v: %v", test.src, test.dst, err)
		}
	}

}
