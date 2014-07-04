package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"
	"syscall"
)

func main() {
	var dependencies []StaticLibrary
	parallels := 0
	done := make(chan bool)

	// set parallel numbers
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Parse flags
	version := flag.String("v", NGINX_VERSION, "nginx version")
	nginxConfPath := flag.String("c", "", "configuration file for building nginx")
	modulesConfPath := flag.String("m", "", "configuration file for 3rd party modules")
	workParentDir := flag.String("d", "", "working directory")
	jobs := flag.Int("j", runtime.NumCPU(), "jobs to build nginx")
	verbose := flag.Bool("verbose", false, "verbose mode")
	pcreStatic := flag.Bool("pcre", false, "embedded PCRE statically")
	pcreVersion := flag.String("pcreversion", PCRE_VERSION, "PCRE version")
	openSSLStatic := flag.Bool("openssl", false, "embedded PCRE statically")
	openSSLVersion := flag.String("opensslversion", OPENSSL_VERSION, "OpenSSL version")
	zlibStatic := flag.Bool("zlib", false, "embedded zlib statically")
	zlibVersion := flag.String("zlibversion", ZLIB_VERSION, "zlib version")
	clear := flag.Bool("clear", false, "remove entries in working directory")
	versionb := flag.Bool("version", false, "print nginx-build versions")
	versions := flag.Bool("versions", false, "print nginx versions")
	flag.Parse()

	if *versionb {
		printNginxBuildVersion()
		os.Exit(0)
	}

	if *versions {
		printNginxVersions()
		os.Exit(0)
	}

	printFirstMsg()

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

	if len(*workParentDir) == 0 {
		log.Fatal("set working directory with -d")
	}

	if !fileExists(*workParentDir) {
		log.Fatalf("working directory(%s) does not exist.", *workParentDir)
	}

	workDir := *workParentDir + "/" + *version
	if *clear {
		err := os.RemoveAll(workDir)
		if err != nil {
			// workaround for a restriction of os.RemoveAll
			// os.RemoveAll call fd.Readdirnames(100).
			// So os.RemoveAll does not always remove all entries.
			// Some 3rd-party module(e.g. lua-nginx-module) tumbles this restriction.
			if fileExists(workDir) {
				err := os.RemoveAll(workDir)
				if err != nil {
					log.Fatal(err.Error())
				}
			} else {
				log.Fatal(err.Error())
			}
		}
	}

	if !fileExists(workDir) {
		err := os.Mkdir(workDir, 0755)
		if err != nil {
			log.Fatalf("Failed to create working directory(%s) does not exist.", workDir)
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
		go downloadAndExtractParallel(&pcreBuilder, done)
	}

	if *openSSLStatic {
		parallels++
		go downloadAndExtractParallel(&openSSLBuilder, done)
	}

	if *zlibStatic {
		parallels++
		go downloadAndExtractParallel(&zlibBuilder, done)
	}

	parallels++
	go downloadAndExtractParallel(&nginxBuilder, done)

	if len(modules3rd) > 0 {
		parallels += len(modules3rd)
		for _, m := range modules3rd {
			go downloadAndExtractModule3rdParallel(m, done)
		}

	}

	// wait until all downloading processes by goroutine finish
	for i := 0; i < parallels; i++ {
		<-done
	}

	if len(modules3rd) > 0 {
		for _, m := range modules3rd {
			provideModule3rd(&m)
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

	log.Printf("Generate configure script for %s.....", nginxBuilder.sourcePath())

	if *pcreStatic && strings.Contains(nginxConf, pcreBuilder.option()+"=") {
		log.Printf("[warn]Using '%s' is discouraged. Instead give '-pcre' and '-pcreversion' to 'nginx-build'", pcreBuilder.option())
	}

	if *openSSLStatic && strings.Contains(nginxConf, openSSLBuilder.option()+"=") {
		log.Printf("[warn]Using '%s' is discouraged. Instead give '-openssl' and '-opensslversion' to 'nginx-build'", openSSLBuilder.option())

	}

	if *zlibStatic && strings.Contains(nginxConf, zlibBuilder.option()+"=") {
		log.Printf("[warn]Using '%s' is discouraged. Instead give '-zlib' and '-zlibversion' to 'nginx-build'", zlibBuilder.option())
	}

	if strings.Contains(nginxConf, "--add-module=") {
		log.Println("[warn]Using '--add-module' is discouraged. Instead give ini-file with '-m' to 'nginx-build'")
	}

	configureScript := configureGen(nginxConf, modules3rd, dependencies)
	err = ioutil.WriteFile("./nginx-configure", []byte(configureScript), 0655)
	if err != nil {
		log.Fatalf("Failed to generate configure script for %s", nginxBuilder.sourcePath())
	}
	log.Printf("Configure %s.....", nginxBuilder.sourcePath())
	err = configureNginx()
	if err != nil {
		log.Fatalf("Failed to configure %s", nginxBuilder.sourcePath())
	}

	log.Printf("Build %s.....", nginxBuilder.sourcePath())
	if *openSSLStatic {
		// This is a workaround for protecting a failure of building nginx with OpenSSL.
		// Unfortunately build of nginx with static OpenSSL fails by multi-CPUs.
		*jobs = 1
	}
	err = buildNginx(*jobs)
	if err != nil {
		fmt.Println(err.Error())
		log.Fatalf("Failed to build %s", nginxBuilder.sourcePath())
	}

	printLastMsg(workDir, nginxBuilder.sourcePath())

	// cd rootDir
	os.Chdir(rootDir)
}
