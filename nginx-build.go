package main

import (
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

func printLastMsg(workDir, srcDir string) {
	log.Println("Complete building nginx!")

	lastMsgFormat := `Enter the following command for install nginx.

$ cd %s/%s
$ sudo make install
`
	log.Println(fmt.Sprintf(lastMsgFormat, workDir, srcDir))
}

func main() {
	version := flag.String("v", NGINX_VERSION, "nginx version")
	confPath := flag.String("c", "", "configuration file for building nginx")
	modulesConfPath := flag.String("m", "", "configuration file for 3rd party modules")
	workParentDir := flag.String("d", "", "working directory")
	verbose := flag.Bool("verbose", false, "verbose mode")
	pcreStatic := flag.Bool("pcre", false, "embedded PCRE statically")
	pcreVersion := flag.String("pcreversion", PCRE_VERSION, "PCRE version")
	openSSLStatic := flag.Bool("openssl", false, "embedded PCRE statically")
	openSSLVersion := flag.String("opensslversion", OPENSSL_VERSION, "OpenSSL version")
	clear := flag.Bool("clear", false, "remove entries in working directory")
	jobs := flag.Int("j", runtime.NumCPU(), "number of jobs for buiding nginx")
	flag.Parse()

	var modulesConf *config.Config
	var modules3rd []Module3rd
	conf := ""
	Verboseenabled = *verbose

	nginxBuilder := MakeBuilder(COMPONENT_NGINX, *version)
	pcreBuilder := MakeBuilder(COMPONENT_PCRE, *pcreVersion)
	openSSLBuilder := MakeBuilder(COMPONENT_OPENSSL, *openSSLVersion)

	// change default umask
	_ = syscall.Umask(0)

	if *version == "" {
		log.Println("[warn]nginx version is not set.")
		log.Printf("[warn]nginx-build use %s.\n", NGINX_VERSION)
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
		modules3rd = LoadModules3rd(modulesConf)
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
		_, err = os.Stat(pcreBuilder.SourcePath())
		if err != nil {
			log.Println("Download PCRE.....")
			err := pcreBuilder.Download()
			if err != nil {
				log.Fatal("Failed to download PCRE")
			}
			err = ExtractArchive(pcreBuilder.ArchivePath())
			if err != nil {
				log.Fatal("Failed to extract nginx")
			}
		} else {
			log.Println(pcreBuilder.SourcePath(), "already exists.")
		}
	}

	if *openSSLStatic {
		_, err = os.Stat(openSSLBuilder.SourcePath())
		if err != nil {
			log.Println("Download OpenSSL.....")
			err := openSSLBuilder.Download()
			if err != nil {
				log.Fatal("Failed to download OpenSSL")
			}
			err = ExtractArchive(openSSLBuilder.ArchivePath())
			if err != nil {
				log.Fatal("Failed to extract nginx")
			}
		} else {
			log.Println(openSSLBuilder.SourcePath(), "already exists.")
		}
	}

	_, err = os.Stat(nginxBuilder.SourcePath())
	if err != nil {
		log.Println("Download nginx.....")
		err := nginxBuilder.Download()
		if err != nil {
			log.Fatal("Failed to download nginx")
		}
		log.Println("Extract nginx.....")
		err = ExtractArchive(nginxBuilder.ArchivePath())
		if err != nil {
			log.Fatal("Failed to extract nginx")
		}
	} else {
		log.Println(nginxBuilder.SourcePath(), "already exists.")
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
			err = nginxBuilder.DownloadModule3rd(modules3rd[i])
			if err != nil {
				log.Fatal("Failed to download ", modules3rd[i].Name)
			}

			if modules3rd[i].Rev != "" {
				dir := saveCurrentDir()
				os.Chdir(modules3rd[i].Name)
				SwitchRev(modules3rd[i].Rev)
				os.Chdir(dir)
			}
		}
	}

	// cd workDir/nginx-${version}
	os.Chdir(nginxBuilder.SourcePath())

	log.Println("Configure nginx.....")
	err = nginxBuilder.ConfigureGen(conf, modules3rd, *pcreStatic, *pcreVersion, *openSSLStatic, *openSSLVersion)
	if err != nil {
		log.Fatal("Failed to generate configure script for nginx")
	}
	err = Configure()
	if err != nil {
		log.Fatal("Failed to configure nginx")
	}

	log.Println("Building nginx.....")
	err = Make(*jobs)
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal("Failed to build nginx")
	}

	// cd rootDir
	os.Chdir(rootDir)

	printLastMsg(workDir, nginxBuilder.SourcePath())
}
