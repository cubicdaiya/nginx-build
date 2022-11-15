package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func untargz(dst, src string) error {

	r, err := os.Open(src)
	if err != nil {
		return err
	}
	defer r.Close()

	gr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gr.Close()

	tr := tar.NewReader(gr)

	for {
		h, err := tr.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		ent := filepath.Join(dst, h.Name)

		if h.FileInfo().Mode().IsDir() {
			if _, err := os.Stat(ent); err != nil {
				if err := os.Mkdir(ent, os.FileMode(h.Mode)); err != nil {
					return err
				}
			}
		} else if h.FileInfo().Mode().IsRegular() {
			f, err := os.OpenFile(ent, os.O_CREATE|os.O_RDWR, os.FileMode(h.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}
			// not use defer in loop
			f.Close()
		} else {
			return fmt.Errorf("The archive contains an entry except regular file and directory: %s", ent)
		}
	}

	return nil
}

func extractArchive(dst, src string) error {
	return untargz(dst, src)
}
