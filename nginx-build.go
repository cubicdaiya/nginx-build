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

	// Parse flags
	version := flag.String("v", NGINX_VERSION, "nginx version")
	nginxConfPath := flag.String("c", "", "configuration file for building nginx")
	modulesConfPath := flag.String("m", "", "configuration file for 3rd party modules")
	workParentDir := flag.String("d", "", "working directory")
	verbose := flag.Bool("verbose", false, "verbose mode")
	pcreStatic := flag.Bool("pcre", false, "embedded PCRE statically")
	pcreVersion := flag.String("pcreversion", PCRE_VERSION, "PCRE version")
	openSSLStatic := flag.Bool("openssl", false, "embedded PCRE statically")
	openSSLVersion := flag.String("zlibversion", OPENSSL_VERSION, "OpenSSL version")
	zlibStatic := flag.Bool("zlib", false, "embedded zlib statically")
	zlibVersion := flag.String("opensslversion", ZLIB_VERSION, "zlib version")
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

	if *workParentDir == "" {
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
		err = downloadAndExtract(&pcreBuilder)
		if err != nil {
			log.Fatal(err.Error())
		}
		dependencies = append(dependencies,
			StaticLibrary{
				Name:    pcreBuilder.name(),
				Version: *pcreVersion,
				Option:  "--with-pcre"})
	}

	if *openSSLStatic {
		err := downloadAndExtract(&openSSLBuilder)
		if err != nil {
			log.Fatal(err.Error())
		}
		dependencies = append(dependencies,
			StaticLibrary{
				Name:    openSSLBuilder.name(),
				Version: *openSSLVersion,
				Option:  "--with-openssl"})
	}

	if *zlibStatic {
		err := downloadAndExtract(&zlibBuilder)
		if err != nil {
			log.Fatal(err.Error())
		}
		dependencies = append(dependencies,
			StaticLibrary{
				Name:    zlibBuilder.name(),
				Version: *zlibVersion,
				Option:  "--with-zlib"})
	}

	err = downloadAndExtract(&nginxBuilder)
	if err != nil {
		log.Fatal(err.Error())
	}
	if len(modules3rd) > 0 {
		log.Println("Download 3rd Party Modules.....")
		for _, m := range modules3rd {
			if fileExists(m.Name) {
				log.Printf("%s already exists.", m.Name)
				continue
			}
			log.Printf("Download %s.....", m.Name)
			err = downloadModule3rd(m)
			if err != nil {
				log.Fatalf("Failed to download %s", m.Name)
			}

			if m.Rev != "" {
				dir := saveCurrentDir()
				os.Chdir(m.Name)
				switchRev(m.Rev)
				prevShell(m.PrevSh)
				os.Chdir(dir)
			}
		}
	}

	// cd workDir/nginx-${version}
	os.Chdir(nginxBuilder.sourcePath())

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
	err = make(*jobs)
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal("Failed to build nginx")
	}

	// cd rootDir
	os.Chdir(rootDir)

	printLastMsg(workDir, nginxBuilder.sourcePath())
}
