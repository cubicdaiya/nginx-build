package main

import (
	"flag"
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
		} else {
			v.Value = flag.String(k, "", v.Desc)
			argsString[k] = v
		}
	}

	for k, v := range argsBool {
		v.Enabled = flag.Bool(k, false, v.Desc)
		argsBool[k] = v
	}

	flag.Parse()

	// Allow multiple flags for `--add-module`
	{
		tmp := argsString["add-module"]
		tmp_ := multiflag.String()
		tmp.Value = &tmp_
		argsString["add-module"] = tmp
	}

	configureOptions.Values = argsString
	configureOptions.Bools = argsBool

	if *versionPrint {
		printNginxBuildVersion()
		os.Exit(0)
	}

	if *versionsPrint {
		printNginxVersions()
		os.Exit(0)
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
		log.Fatal(err.Error())
	}
	nginxConfigure = normalizeConfigure(nginxConfigure)

	modules3rd, err := loadModules3rdFile(*modulesConfPath)
	if err != nil {
		log.Fatal(err.Error())
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
		workDir = *workParentDir + "/" + *openRestyVersion
	} else if *tengine {
		workDir = *workParentDir + "/" + *tengineVersion
	} else {
		workDir = *workParentDir + "/" + *version
	}
	if *clear {
		err := clearWorkDir(workDir)
		if err != nil {
			log.Fatal(err.Error())
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

	if *pcreStatic && pcreBuilder.isIncludeWithOption(nginxConfigure) {
		log.Printf(pcreBuilder.warnMsgWithLibrary())
	}

	if *openSSLStatic && openSSLBuilder.isIncludeWithOption(nginxConfigure) {
		log.Printf(openSSLBuilder.warnMsgWithLibrary())

	}

	if *zlibStatic && zlibBuilder.isIncludeWithOption(nginxConfigure) {
		log.Printf(zlibBuilder.warnMsgWithLibrary())
	}

	if strings.Contains(nginxConfigure, "--add-module=") {
		log.Println("[warn]Using '--add-module' is discouraged. Instead give ini-file with '-m' to 'nginx-build'")
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
		fatalLog(err, "nginx-configure.log")
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
		fatalLog(err, "nginx-build.log")
	}

	printLastMsg(workDir, nginxBuilder.sourcePath(), *openResty, *configureOnly)

	// cd rootDir
	os.Chdir(rootDir)
}
