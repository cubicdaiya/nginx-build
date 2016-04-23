package main

import (
	"bufio"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"github.com/cubicdaiya/nginx-build/builder"
	"github.com/cubicdaiya/nginx-build/openssl"
)

const NGINX_BUILD_VERSION = "0.9.1"

func buildNginx(jobs int) error {
	args := []string{"make", "-j", strconv.Itoa(jobs)}
	if VerboseEnabled {
		return runCommand(args)
	}

	f, err := os.Create("nginx-build.log")
	if err != nil {
		return runCommand(args)
	}
	defer f.Close()

	cmd, err := makeCmd(args)
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(f)
	cmd.Stderr = writer
	defer writer.Flush()

	return cmd.Run()
}

func main() {
	var dependencies []builder.StaticLibrary
	wg := new(sync.WaitGroup)

	// Parse flags
	verbose := flag.Bool("verbose", false, "verbose mode")
	pcreStatic := flag.Bool("pcre", false, "embedded PCRE statically")
	openSSLStatic := flag.Bool("openssl", false, "embedded OpenSSL statically")
	zlibStatic := flag.Bool("zlib", false, "embedded zlib statically")
	clear := flag.Bool("clear", false, "remove entries in working directory")
	versionPrint := flag.Bool("version", false, "print nginx-build versions")
	versionsPrint := flag.Bool("versions", false, "print nginx versions")
	openResty := flag.Bool("openresty", false, "download openresty instead of nginx")
	tengine := flag.Bool("tengine", false, "download tengine instead of nginx")
	configureOnly := flag.Bool("configureonly", false, "configuring nginx only not building")
	jobs := flag.Int("j", runtime.NumCPU(), "jobs to build nginx")
	version := flag.String("v", builder.NGINX_VERSION, "nginx version")
	nginxConfigurePath := flag.String("c", "", "configuration file for building nginx")
	modulesConfPath := flag.String("m", "", "configuration file for 3rd party modules")
	workParentDir := flag.String("d", "", "working directory")
	pcreVersion := flag.String("pcreversion", builder.PCRE_VERSION, "PCRE version")
	openSSLVersion := flag.String("opensslversion", builder.OPENSSL_VERSION, "OpenSSL version")
	zlibVersion := flag.String("zlibversion", builder.ZLIB_VERSION, "zlib version")
	openRestyVersion := flag.String("openrestyversion", builder.OPENRESTY_VERSION, "openresty version")
	tengineVersion := flag.String("tengineversion", builder.TENGINE_VERSION, "tengine version")

	// fake flag for --with-xxx=dynamic
	for i, arg := range os.Args {
		if strings.Contains(arg, "with-http_xslt_module=dynamic") {
			os.Args[i] = "--with-http_xslt_module_dynamic"
		}
		if strings.Contains(arg, "with-http_image_filter_module=dynamic") {
			os.Args[i] = "--with-http_image_filter_module_dynamic"
		}
		if strings.Contains(arg, "with-http_geoip_module=dynamic") {
			os.Args[i] = "--with-http_geoip_module_dynamic"
		}
		if strings.Contains(arg, "with-http_perl_module=dynamic") {
			os.Args[i] = "--with-http_perl_module_dynamic"
		}
		if strings.Contains(arg, "with-mail=dynamic") {
			os.Args[i] = "--with-mail_dynamic"
		}
		if strings.Contains(arg, "with-stream=dynamic") {
			os.Args[i] = "--with-stream_dynamic"
		}
	}

	var (
		configureOptions ConfigureOptions
		multiflag        StringFlag
		multiflagDynamic StringFlag
	)

	argsString := makeArgsString()
	for k, v := range argsString {
		if k == "add-module" {
			flag.Var(&multiflag, k, v.Desc)
		} else if k == "add-dynamic-module" {
			flag.Var(&multiflagDynamic, k, v.Desc)
		} else {
			v.Value = flag.String(k, "", v.Desc)
			argsString[k] = v
		}
	}

	argsBool := makeArgsBool()
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
		tmp_ := multiflagDynamic.String()
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

	var nginxBuilder builder.Builder
	if *openResty && *tengine {
		log.Fatal("select one between '-openresty' and '-tengine'.")
	}
	if *openResty {
		nginxBuilder = builder.MakeBuilder(builder.COMPONENT_OPENRESTY, *openRestyVersion)
	} else if *tengine {
		nginxBuilder = builder.MakeBuilder(builder.COMPONENT_TENGINE, *tengineVersion)
	} else {
		nginxBuilder = builder.MakeBuilder(builder.COMPONENT_NGINX, *version)
	}
	pcreBuilder := builder.MakeBuilder(builder.COMPONENT_PCRE, *pcreVersion)
	openSSLBuilder := builder.MakeBuilder(builder.COMPONENT_OPENSSL, *openSSLVersion)
	zlibBuilder := builder.MakeBuilder(builder.COMPONENT_ZLIB, *zlibVersion)

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
	os.Chdir(nginxBuilder.SourcePath())

	if *pcreStatic {
		dependencies = append(dependencies, builder.MakeStaticLibrary(&pcreBuilder))
	}

	if *openSSLStatic {
		dependencies = append(dependencies, builder.MakeStaticLibrary(&openSSLBuilder))
	}

	if *zlibStatic {
		dependencies = append(dependencies, builder.MakeStaticLibrary(&zlibBuilder))
	}

	log.Printf("Generate configure script for %s.....", nginxBuilder.SourcePath())

	if *pcreStatic && pcreBuilder.IsIncludeWithOption(nginxConfigure) {
		log.Println(pcreBuilder.WarnMsgWithLibrary())
	}

	if *openSSLStatic && openSSLBuilder.IsIncludeWithOption(nginxConfigure) {
		log.Println(openSSLBuilder.WarnMsgWithLibrary())

	}

	if *zlibStatic && zlibBuilder.IsIncludeWithOption(nginxConfigure) {
		log.Println(zlibBuilder.WarnMsgWithLibrary())
	}

	configureScript := configureGen(nginxConfigure, modules3rd, dependencies, configureOptions, rootDir)

	err = ioutil.WriteFile("./nginx-configure", []byte(configureScript), 0655)
	if err != nil {
		log.Fatalf("Failed to generate configure script for %s", nginxBuilder.SourcePath())
	}

	log.Printf("Configure %s.....", nginxBuilder.SourcePath())

	err = configureNginx()
	if err != nil {
		log.Printf("Failed to configure %s\n", nginxBuilder.SourcePath())
		printFatalMsg(err, "nginx-configure.log")
	}

	if *configureOnly {
		printLastMsg(workDir, nginxBuilder.SourcePath(), *openResty, *configureOnly)
		return
	}

	log.Printf("Build %s.....", nginxBuilder.SourcePath())

	if *openSSLStatic {
		// Workarounds for protecting a failure of building nginx with static-linked OpenSSL.

		// The build of old OpenSSL fails when multi-CPUs are used.
		// It is fixed in the commit below.
		//
		// https://github.com/openssl/openssl/commit/8e6bb99979b95ee8b878e22e043ceb78d79c32a1
		//
		// And backported to 1.0.1p and 1.0.2d.
		if !openssl.ParallelBuildAvailable(*openSSLVersion) {
			*jobs = 1
		}

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
		log.Printf("Failed to build %s\n", nginxBuilder.SourcePath())
		printFatalMsg(err, "nginx-build.log")
	}

	printLastMsg(workDir, nginxBuilder.SourcePath(), *openResty, *configureOnly)

	// cd rootDir
	os.Chdir(rootDir)
}
