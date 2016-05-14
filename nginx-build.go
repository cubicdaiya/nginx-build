package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"syscall"

	"github.com/cubicdaiya/nginx-build/builder"
	"github.com/cubicdaiya/nginx-build/command"
	"github.com/cubicdaiya/nginx-build/configure"
	"github.com/cubicdaiya/nginx-build/module3rd"
	"github.com/cubicdaiya/nginx-build/openssl"
	"github.com/cubicdaiya/nginx-build/util"
)

const (
	NginxBuildVersion = "0.9.3"
)

var (
	nginxBuildOptions Options
)

func init() {
	nginxBuildOptions = makeNginxBuildOptions()
}

func main() {
	// Parse flags
	for k, v := range nginxBuildOptions.Bools {
		v.Enabled = flag.Bool(k, false, v.Desc)
		nginxBuildOptions.Bools[k] = v
	}
	for k, v := range nginxBuildOptions.Values {
		v.Value = flag.String(k, v.Default, v.Desc)
		nginxBuildOptions.Values[k] = v
	}
	for k, v := range nginxBuildOptions.Numbers {
		v.Value = flag.Int(k, v.Default, v.Desc)
		nginxBuildOptions.Numbers[k] = v
	}

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
		configureOptions configure.Options
		multiflag        StringFlag
		multiflagDynamic StringFlag
	)

	argsString := configure.MakeArgsString()
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

	argsBool := configure.MakeArgsBool()
	for k, v := range argsBool {
		v.Enabled = flag.Bool(k, false, v.Desc)
		argsBool[k] = v
	}

	flag.CommandLine.SetOutput(os.Stdout)
	// The output of original flag.Usage() is too long
	defaultUsage := flag.Usage
	flag.Usage = usage
	flag.Parse()

	jobs := nginxBuildOptions.Numbers["j"].Value

	verbose := nginxBuildOptions.Bools["verbose"].Enabled
	pcreStatic := nginxBuildOptions.Bools["pcre"].Enabled
	openSSLStatic := nginxBuildOptions.Bools["openssl"].Enabled
	zlibStatic := nginxBuildOptions.Bools["zlib"].Enabled
	clear := nginxBuildOptions.Bools["clear"].Enabled
	versionPrint := nginxBuildOptions.Bools["version"].Enabled
	versionsPrint := nginxBuildOptions.Bools["versions"].Enabled
	openResty := nginxBuildOptions.Bools["openresty"].Enabled
	tengine := nginxBuildOptions.Bools["tengine"].Enabled
	configureOnly := nginxBuildOptions.Bools["configureonly"].Enabled
	helpAll := nginxBuildOptions.Bools["help-all"].Enabled

	version := nginxBuildOptions.Values["v"].Value
	nginxConfigurePath := nginxBuildOptions.Values["c"].Value
	modulesConfPath := nginxBuildOptions.Values["m"].Value
	workParentDir := nginxBuildOptions.Values["d"].Value
	pcreVersion := nginxBuildOptions.Values["pcreversion"].Value
	openSSLVersion := nginxBuildOptions.Values["opensslversion"].Value
	zlibVersion := nginxBuildOptions.Values["zlibversion"].Value
	openRestyVersion := nginxBuildOptions.Values["openrestyversion"].Value
	tengineVersion := nginxBuildOptions.Values["tengineversion"].Value

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

	if *helpAll {
		defaultUsage()
		return
	}

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
	command.VerboseEnabled = *verbose

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

	nginxConfigure, err := util.FileGetContents(*nginxConfigurePath)
	if err != nil {
		log.Fatal(err)
	}
	nginxConfigure = configure.Normalize(nginxConfigure)

	modules3rd, err := module3rd.Load(*modulesConfPath)
	if err != nil {
		log.Fatal(err)
	}

	if len(*workParentDir) == 0 {
		log.Fatal("set working directory with -d")
	}

	if !util.FileExists(*workParentDir) {
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
		err := util.ClearWorkDir(workDir)
		if err != nil {
			log.Fatal(err)
		}
	}

	if !util.FileExists(workDir) {
		err := os.MkdirAll(workDir, 0755)
		if err != nil {
			log.Fatalf("Failed to create working directory(%s) does not exist.", workDir)
		}
	}

	rootDir := util.SaveCurrentDir()
	// cd workDir
	err = os.Chdir(workDir)
	if err != nil {
		log.Fatal(err)
	}

	wg := new(sync.WaitGroup)
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
			go module3rd.DownloadAndExtractParallel(m, wg)
		}

	}

	// wait until all downloading processes by goroutine finish
	wg.Wait()

	if len(modules3rd) > 0 {
		for _, m := range modules3rd {
			if err := module3rd.Provide(&m); err != nil {
				log.Fatal(err)
			}
		}
	}

	// cd workDir/nginx-${version}
	os.Chdir(nginxBuilder.SourcePath())

	var dependencies []builder.StaticLibrary
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

	configureScript := configure.Generate(nginxConfigure, modules3rd, dependencies, configureOptions, rootDir, *openResty, *jobs)

	err = ioutil.WriteFile("./nginx-configure", []byte(configureScript), 0655)
	if err != nil {
		log.Fatalf("Failed to generate configure script for %s", nginxBuilder.SourcePath())
	}

	log.Printf("Configure %s.....", nginxBuilder.SourcePath())

	err = configure.Run()
	if err != nil {
		log.Printf("Failed to configure %s\n", nginxBuilder.SourcePath())
		util.PrintFatalMsg(err, "nginx-configure.log")
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

	err = builder.BuildNginx(*jobs)
	if err != nil {
		log.Printf("Failed to build %s\n", nginxBuilder.SourcePath())
		util.PrintFatalMsg(err, "nginx-build.log")
	}

	printLastMsg(workDir, nginxBuilder.SourcePath(), *openResty, *configureOnly)

	// cd rootDir
	os.Chdir(rootDir)
}
