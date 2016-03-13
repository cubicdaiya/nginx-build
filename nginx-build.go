package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"sync"
	"syscall"
)

func main() {
	var dependencies []StaticLibrary
	wg := new(sync.WaitGroup)

	// Parse flags
	version := flag.String("v", NGINX_VERSION, "nginx version")
	nginxConfigurePath := flag.String("c", "", "configuration file for building nginx")
	modulesConfPath := flag.String("m", "", "configuration file for 3rd party modules")
	workParentDir := flag.String("d", "", "working directory")
	jobs := flag.Int("j", runtime.NumCPU(), "jobs to build nginx")
	verbose := flag.Bool("verbose", false, "verbose mode")
	pcreStatic := flag.Bool("pcre", false, "embedded PCRE statically")
	pcreVersion := flag.String("pcreversion", PCRE_VERSION, "PCRE version")
	openSSLStatic := flag.Bool("openssl", false, "embedded OpenSSL statically")
	openSSLVersion := flag.String("opensslversion", OPENSSL_VERSION, "OpenSSL version")
	zlibStatic := flag.Bool("zlib", false, "embedded zlib statically")
	zlibVersion := flag.String("zlibversion", ZLIB_VERSION, "zlib version")
	clear := flag.Bool("clear", false, "remove entries in working directory")
	versionPrint := flag.Bool("version", false, "print nginx-build versions")
	versionsPrint := flag.Bool("versions", false, "print nginx versions")
	openResty := flag.Bool("openresty", false, "download openresty instead of nginx")
	tengine := flag.Bool("tengine", false, "download tengine instead of nginx")
	openRestyVersion := flag.String("openrestyversion", OPENRESTY_VERSION, "openresty version")
	tengineVersion := flag.String("tengineversion", TENGINE_VERSION, "tengine version")
	configureOnly := flag.Bool("configureonly", false, "configuring nginx only not building")

	var configureOptions ConfigureOptions
	argsString := makeArgsString()
	argsBool := makeArgsBool()
	var multiflag StringFlag

	for k, v := range argsString {
		if k == "add-module" {
			flag.Var(&multiflag, k, v.Desc)
		} else if k == "add-dynamic-module" {
			flag.Var(&multiflag, k, v.Desc)
		} else {
			v.Value = flag.String(k, "", v.Desc)
			argsString[k] = v
		}
	}

	for k, v := range argsBool {
		v.Enabled = flag.Bool(k, false, v.Desc)
		argsBool[k] = v
	}

	flag.CommandLine.SetOutput(os.Stdout)
	flag.Parse()

	// Allow multiple flags for `--add-module`
	{
		tmp := argsString["add-module"]
		tmp_ := multiflag.String()
		tmp.Value = &tmp_
		argsString["add-module"] = tmp
	}

	// Allow multiple flags for `--add-dynamic-module`
	{
		tmp := argsString["add-dynamic-module"]
		tmp_ := multiflag.String()
		tmp.Value = &tmp_
		argsString["add-dynamic-module"] = tmp
	}

	configureOptions.Values = argsString
	configureOptions.Bools = argsBool

	if *versionPrint {
		printNginxBuildVersion()
		return
	}

	if *versionsPrint {
		printNginxVersions()
		return
	}

	printFirstMsg()

	// set verbose mode
	VerboseEnabled = *verbose

	var nginxBuilder Builder
	if *openResty && *tengine {
		log.Fatal("select one between '-openresty' and '-tengine'.")
	}
	if *openResty {
		nginxBuilder = makeBuilder(COMPONENT_OPENRESTY, *openRestyVersion)
	} else if *tengine {
		nginxBuilder = makeBuilder(COMPONENT_TENGINE, *tengineVersion)
	} else {
		nginxBuilder = makeBuilder(COMPONENT_NGINX, *version)
	}
	pcreBuilder := makeBuilder(COMPONENT_PCRE, *pcreVersion)
	openSSLBuilder := makeBuilder(COMPONENT_OPENSSL, *openSSLVersion)
	zlibBuilder := makeBuilder(COMPONENT_ZLIB, *zlibVersion)

	// change default umask
	_ = syscall.Umask(0)

	versionCheck(*version)

	nginxConfigure, err := fileGetContents(*nginxConfigurePath)
	if err != nil {
		log.Fatal(err)
	}
	nginxConfigure = normalizeConfigure(nginxConfigure)

	modules3rd, err := loadModules3rdFile(*modulesConfPath)
	if err != nil {
		log.Fatal(err)
	}

	if len(*workParentDir) == 0 {
		log.Fatal("set working directory with -d")
	}

	if !fileExists(*workParentDir) {
		err := os.Mkdir(*workParentDir, 0755)
		if err != nil {
			log.Fatalf("Failed to create working directory(%s) does not exist.", *workParentDir)
		}
	}

	var workDir string
	if *openResty {
		workDir = *workParentDir + "/openresty/" + *openRestyVersion
	} else if *tengine {
		workDir = *workParentDir + "/tengine/" + *tengineVersion
	} else {
		workDir = *workParentDir + "/nginx/" + *version
	}
	if *clear {
		err := clearWorkDir(workDir)
		if err != nil {
			log.Fatal(err)
		}
	}

	if !fileExists(workDir) {
		err := os.MkdirAll(workDir, 0755)
		if err != nil {
			log.Fatalf("Failed to create working directory(%s) does not exist.", workDir)
		}
	}

	rootDir := saveCurrentDir()
	// cd workDir
	err = os.Chdir(workDir)
	if err != nil {
		log.Fatal(err)
	}

	if *pcreStatic {
		wg.Add(1)
		go downloadAndExtractParallel(&pcreBuilder, wg)
	}

	if *openSSLStatic {
		wg.Add(1)
		go downloadAndExtractParallel(&openSSLBuilder, wg)
	}

	if *zlibStatic {
		wg.Add(1)
		go downloadAndExtractParallel(&zlibBuilder, wg)
	}

	wg.Add(1)
	go downloadAndExtractParallel(&nginxBuilder, wg)

	if len(modules3rd) > 0 {
		wg.Add(len(modules3rd))
		for _, m := range modules3rd {
			go downloadAndExtractModule3rdParallel(m, wg)
		}

	}

	// wait until all downloading processes by goroutine finish
	wg.Wait()

	if len(modules3rd) > 0 {
		for _, m := range modules3rd {
			if err := provideModule3rd(&m); err != nil {
				log.Fatal(err)
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

	log.Printf("Generate configure script for %s.....", nginxBuilder.sourcePath())

	if *pcreStatic && pcreBuilder.isIncludeWithOption(nginxConfigure) {
		log.Println(pcreBuilder.warnMsgWithLibrary())
	}

	if *openSSLStatic && openSSLBuilder.isIncludeWithOption(nginxConfigure) {
		log.Println(openSSLBuilder.warnMsgWithLibrary())

	}

	if *zlibStatic && zlibBuilder.isIncludeWithOption(nginxConfigure) {
		log.Println(zlibBuilder.warnMsgWithLibrary())
	}

	configureScript := configureGen(nginxConfigure, modules3rd, dependencies, configureOptions, rootDir)

	err = ioutil.WriteFile("./nginx-configure", []byte(configureScript), 0655)
	if err != nil {
		log.Fatalf("Failed to generate configure script for %s", nginxBuilder.sourcePath())
	}

	log.Printf("Configure %s.....", nginxBuilder.sourcePath())

	err = configureNginx()
	if err != nil {
		log.Printf("Failed to configure %s\n", nginxBuilder.sourcePath())
		printFatalMsg(err, "nginx-configure.log")
	}

	if *configureOnly {
		printLastMsg(workDir, nginxBuilder.sourcePath(), *openResty, *configureOnly)
		return
	}

	log.Printf("Build %s.....", nginxBuilder.sourcePath())

	if *openSSLStatic {
		// Workarounds for protecting a failure of building nginx with static-linked OpenSSL.

		// Unfortunately a build of OpenSSL fails when multi-CPUs are used.
		*jobs = 1

		// Sometimes machine hardware name('uname -m') is different
		// from machine processor architecture name('uname -p') on Mac.
		// Specifically, `uname -p` is 'i386' and `uname -m` is 'x86_64'.
		// In this case, a build of OpenSSL fails.
		// So it needs to convince OpenSSL with KERNEL_BITS.
		if runtime.GOOS == "darwin" && runtime.GOARCH == "amd64" {
			os.Setenv("KERNEL_BITS", "64")
		}
	}

	err = buildNginx(*jobs)
	if err != nil {
		log.Printf("Failed to build %s\n", nginxBuilder.sourcePath())
		printFatalMsg(err, "nginx-build.log")
	}

	printLastMsg(workDir, nginxBuilder.sourcePath(), *openResty, *configureOnly)

	// cd rootDir
	os.Chdir(rootDir)
}
