package main

import (
	"./common"
	"./nginx"
	"./openssl"
	"./pcre"
	"flag"
	"fmt"
	"github.com/robfig/config"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"syscall"
)

func saveCurrentDir() string {
	prevDir, _ := filepath.Abs(".")
	return prevDir
}

func restoreCurrentDir(prevDir string) {
	os.Chdir(prevDir)
}

func printLastMsg(workDir, version string) {
	log.Println("Complete building nginx!")

	lastMsgFormat := `Enter the following command for install nginx.

$ cd %s/%s
$ sudo make install
`
	log.Println(fmt.Sprintf(lastMsgFormat, workDir, nginx.SourcePath(version)))
}

func main() {
	version := flag.String("v", nginx.VERSION, "nginx version")
	confPath := flag.String("c", "", "configuration file for building nginx")
	modulesConfPath := flag.String("m", "", "configuration file for 3rd party modules")
	workParentDir := flag.String("d", "", "working directory")
	verbose := flag.Bool("verbose", false, "verbose mode")
	pcreStatic := flag.Bool("pcre", false, "embedded PCRE statically")
	pcreVersion := flag.String("pcreversion", pcre.VERSION, "PCRE version")
	openSSLStatic := flag.Bool("openssl", false, "embedded PCRE statically")
	openSSLVersion := flag.String("opensslversion", openssl.VERSION, "OpenSSL version")
	clear := flag.Bool("clear", false, "remove entries in working directory")
	jobs := flag.Int("j", runtime.NumCPU(), "number of jobs for buiding nginx")
	flag.Parse()

	var modulesConf *config.Config
	var modules3rd []nginx.Module3rd
	conf := ""
	common.Verboseenabled = *verbose

	// change default umask
	_ = syscall.Umask(0)

	if *version == "" {
		log.Println("[warn]nginx version is not set.")
		log.Printf("[warn]nginx-build use %s.\n", nginx.VERSION)
	}

	if *confPath == "" {
		log.Println("[warn]configure option is empty.")
	} else {
		confb, err := ioutil.ReadFile(*confPath)
		if err != nil {
			log.Fatal(fmt.Sprintf("confPath(%s) does not exist.", *confPath))
		}
		conf = string(confb)
	}

	if *modulesConfPath != "" {
		_, err := os.Stat(*modulesConfPath)
		if err != nil {
			log.Fatal(fmt.Sprintf("modulesConfPath(%s) does not exist.", modulesConfPath))
		}

		modulesConf, err = config.ReadDefault(*modulesConfPath)
		if err != nil {
			log.Fatal(err)
		}
		modules3rd = nginx.LoadModules3rd(modulesConf)
	}

	if *workParentDir == "" {
		log.Fatal("set working directory with -d")
	}

	_, err := os.Stat(*workParentDir)
	if err != nil {
		log.Fatal(fmt.Sprintf("working directory(%s) does not exist.", *workParentDir))
	}

	workDir := *workParentDir + "/" + *version
	if *clear {
		err = os.RemoveAll(workDir)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
	_, err = os.Stat(workDir)
	if err != nil {
		err = os.Mkdir(workDir, 0755)
		if err != nil {
			log.Fatal("Failed to create working directory(%s) does not exist.", workDir)
		}
	}

	rootDir := saveCurrentDir()
	// cd workDir
	err = os.Chdir(workDir)
	if err != nil {
		log.Fatal(err.Error())
	}

	if *pcreStatic {
		_, err = os.Stat(pcre.SourcePath(*pcreVersion))
		if err != nil {
			log.Println("Download PCRE.....")
			downloadLink := pcre.DownloadLink(*pcreVersion)
			err := pcre.Download(downloadLink)
			if err != nil {
				log.Fatal("Failed to download PCRE")
			}
			err = nginx.ExtractArchive(pcre.ArchivePath(*pcreVersion))
			if err != nil {
				log.Fatal("Failed to extract nginx")
			}
		} else {
			log.Println(pcre.SourcePath(*pcreVersion), "already exists.")
		}
	}

	if *openSSLStatic {
		_, err = os.Stat(openssl.SourcePath(*openSSLVersion))
		if err != nil {
			log.Println("Download OpenSSL.....")
			downloadLink := openssl.DownloadLink(*openSSLVersion)
			err := openssl.Download(downloadLink)
			if err != nil {
				log.Fatal("Failed to download OpenSSL")
			}
			err = nginx.ExtractArchive(openssl.ArchivePath(*openSSLVersion))
			if err != nil {
				log.Fatal("Failed to extract nginx")
			}
		} else {
			log.Println(openssl.SourcePath(*openSSLVersion), "already exists.")
		}
	}

	_, err = os.Stat(nginx.SourcePath(*version))
	if err != nil {
		log.Println("Download nginx.....")
		downloadLink := nginx.DownloadLink(*version)
		err := nginx.Download(downloadLink)
		if err != nil {
			log.Fatal("Failed to download nginx")
		}
		log.Println("Extract nginx.....")
		err = nginx.ExtractArchive(nginx.ArchivePath(*version))
		if err != nil {
			log.Fatal("Failed to extract nginx")
		}
	} else {
		log.Println(nginx.SourcePath(*version), "already exists.")
	}

	if len(modules3rd) > 0 {
		log.Println("Download 3rd Party Modules.....")
		for i := 0; i < len(modules3rd); i++ {
			_, err := os.Stat(modules3rd[i].Name)
			if err == nil {
				log.Println(modules3rd[i].Name, " already exists.")
				continue
			}
			log.Println(fmt.Sprintf("Download %s.....", modules3rd[i].Name))
			err = nginx.DownloadModule3rd(modules3rd[i])
			if err != nil {
				log.Fatal("Failed to download ", modules3rd[i].Name)
			}

			if modules3rd[i].Rev != "" {
				dir := saveCurrentDir()
				os.Chdir(modules3rd[i].Name)
				nginx.SwitchRev(modules3rd[i].Rev)
				os.Chdir(dir)
			}
		}
	}

	// cd workDir/nginx-${version}
	os.Chdir(nginx.SourcePath(*version))

	log.Println("Configure nginx.....")
	err = nginx.ConfigureGen(conf, modules3rd, *pcreStatic, *pcreVersion, *openSSLStatic, *openSSLVersion)
	if err != nil {
		log.Fatal("Failed to generate configure script for nginx")
	}
	err = nginx.Configure()
	if err != nil {
		log.Fatal("Failed to configure nginx")
	}

	log.Println("Building nginx.....")
	err = nginx.Make(conf, *jobs)
	if err != nil {
		log.Fatal("Failed to build nginx")
	}

	// cd rootDir
	os.Chdir(rootDir)

	printLastMsg(workDir, *version)
}
