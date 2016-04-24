package util

import (
	"bufio"
	"log"
	"os"

	"github.com/cubicdaiya/nginx-build/command"
)

func PrintFatalMsg(err error, path string) {
	if command.VerboseEnabled {
		log.Fatal(err)
	}

	f, err2 := os.Open(path)
	if err2 != nil {
		log.Printf("error-log: %s is not found\n", path)
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		os.Stderr.Write(scanner.Bytes())
		os.Stderr.Write([]byte("\n"))
	}

	log.Fatal(err)
}
