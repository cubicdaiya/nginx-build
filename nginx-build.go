package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"syscall"
)

func main() {
	var dependencies []StaticLibrary
	var parallels int
	done := make(chan bool)

	// set parallel numbers
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Parse flags
	version := flag.String("v", NGINX_VERSION, "nginx version")
	nginxConfPath := flag.String("c", "", "configuration file for building nginx")
	modulesConfPath := flag.String("m", "", "configuration file for 3rd party modules")
	workParentDir := flag.String("d", "", "working directory")
	verbose := flag.Bool("verbose", false, "verbose mode")
	pcreStatic := flag.Bool("pcre", false, "embedded PCRE statically")
	pcreVersion := flag.String("pcreversion", PCRE_VERSION, "PCRE version")
	openSSLStatic := flag.Bool("openssl", false, "embedded PCRE statically")
	openSSLVersion := flag.String("opensslversion", OPENSSL_VERSION, "OpenSSL version")
	zlibStatic := flag.Bool("zlib", false, "embedded zlib statically")
	zlibVersion := flag.String("zlibversion", ZLIB_VERSION, "zlib version")
	clear := flag.Bool("clear", false, "remove entries in working directory")
	jobs := flag.Int("j", runtime.NumCPU(), "number of jobs for buiding nginx")
	versionb := flag.Bool("version", false, "show nginx-build versions")
	versions := flag.Bool("versions", false, "show nginx versions")
	flag.Parse()

	if *versionb {
		showNginxBuildVersion()
		os.Exit(0)
	}

	if *versions {
		showNginxVersions()
		os.Exit(0)
	}

	// set verbose mode
	VerboseEnabled = *verbose

	nginxBuilder := makeBuilder(COMPONENT_NGINX, *version)
	pcreBuilder := makeBuilder(COMPONENT_PCRE, *pcreVersion)
	openSSLBuilder := makeBuilder(COMPONENT_OPENSSL, *openSSLVersion)
	zlibBuilder := makeBuilder(COMPONENT_ZLIB, *zlibVersion)

	// change default umask
	_ = syscall.Umask(0)

	versionCheck(*version)

	nginxConf, err := fileGetContents(*nginxConfPath)
	if err != nil {
		log.Fatal(err.Error())
	}

	modules3rd, err := loadModules3rdFile(*modulesConfPath)
	if err != nil {
		log.Fatal(err.Error())
	}
	parallels += len(modules3rd)

	if len(*workParentDir) == 0 {
		log.Fatal("set working directory with -d")
	}

	if !fileExists(*workParentDir) {
		log.Fatal(fmt.Sprintf("working directory(%s) does not exist.", *workParentDir))
	}

	workDir := *workParentDir + "/" + *version
	if *clear {
		err := os.RemoveAll(workDir)
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	if !fileExists(workDir) {
		err := os.Mkdir(workDir, 0755)
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
		parallels++
		go func(done chan bool) {
			err = downloadAndExtract(&pcreBuilder)
			if err != nil {
				log.Fatal(err.Error())
			}
			done <- true
		}(done)
	}

	if *openSSLStatic {
		parallels++
		go func(done chan bool) {
			err := downloadAndExtract(&openSSLBuilder)
			if err != nil {
				log.Fatal(err.Error())
			}
			done <- true
		}(done)
	}

	if *zlibStatic {
		parallels++
		go func(done chan bool) {
			err := downloadAndExtract(&zlibBuilder)
			if err != nil {
				log.Fatal(err.Error())
			}
			done <- true
		}(done)
	}

	parallels++
	go func(done chan bool) {
		err = downloadAndExtract(&nginxBuilder)
		if err != nil {
			log.Fatal(err.Error())
		}
		done <- true
	}(done)

	if len(modules3rd) > 0 {
		for _, m := range modules3rd {
			go func(done chan bool, m Module3rd) {
				if fileExists(m.Name) {
					log.Printf("%s already exists.", m.Name)
					done <- true
					return
				}
				log.Printf("Download %s.....", m.Name)
				err = downloadModule3rd(m)
				if err != nil {
					log.Println(err.Error())
					log.Fatalf("Failed to download %s", m.Name)
				}
				done <- true
			}(done, m)
		}

	}

	for i := 0; i < parallels; i++ {
		<-done
	}

	if len(modules3rd) > 0 {
		for _, m := range modules3rd {
			if len(m.Rev) > 0 {
				dir := saveCurrentDir()
				os.Chdir(m.Name)
				err := switchRev(m.Rev)
				if err != nil {
					log.Println(err.Error())
				}
				err = prevShell(m.PrevSh)
				if err != nil {
					log.Println(err.Error())
				}
				os.Chdir(dir)
			}
		}
	}

	// cd workDir/nginx-${version}
	os.Chdir(nginxBuilder.sourcePath())

	if *pcreStatic {
		dependencies = append(dependencies, makeStaticLibrary(&pcreBuilder))
	}

	if *openSSLStatic {
		dependencies = append(dependencies, makeStaticLibrary(&openSSLBuilder))
	}

	if *zlibStatic {
		dependencies = append(dependencies, makeStaticLibrary(&zlibBuilder))
	}

	log.Println("Generate configure script for nginx.....")
	err = nginxBuilder.configureGen(nginxConf, modules3rd, dependencies)
	if err != nil {
		log.Fatal("Failed to generate configure script for nginx")
	}
	log.Println("Configure nginx.....")
	err = configure()
	if err != nil {
		log.Fatal("Failed to configure nginx")
	}

	log.Println("Build nginx.....")
	if *openSSLStatic {
		// This is a workaround for protecting a failure of building nginx with OpenSSL.
		// Unfortunately build of nginx with static OpenSSL fails by multi-CPUs.
		*jobs = 1
	}
	err = build(*jobs)
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal("Failed to build nginx")
	}

	// cd rootDir
	os.Chdir(rootDir)

	printLastMsg(workDir, nginxBuilder.sourcePath())
}
