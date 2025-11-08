package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"syscall"

	"github.com/cubicdaiya/nginx-build/builder"
	"github.com/cubicdaiya/nginx-build/command"
	"github.com/cubicdaiya/nginx-build/configure"
	"github.com/cubicdaiya/nginx-build/module3rd"
	"github.com/cubicdaiya/nginx-build/util"
)

const defaultPatchTarget = "nginx"

type patchDirective struct {
	target string
	paths  []string
}

var (
	nginxBuildOptions Options
)

func init() {
	nginxBuildOptions = makeNginxBuildOptions()
}

// fake flag for --with-xxx=dynamic
func overrideUnableParseFlags() {
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
		if strings.Contains(arg, "with-stream_geoip_module=dynamic") {
			os.Args[i] = "--with-stream_geoip_module_dynamic"
		}
	}
}

func main() {
	var (
		multiflagPatch StringFlag
	)

	// Parse flags
	for k, v := range nginxBuildOptions.Bools {
		v.Enabled = flag.Bool(k, false, v.Desc)
		nginxBuildOptions.Bools[k] = v
	}
	for k, v := range nginxBuildOptions.Values {
		if k == "patch" {
			flag.Var(&multiflagPatch, k, v.Desc)
		} else {
			v.Value = flag.String(k, v.Default, v.Desc)
			nginxBuildOptions.Values[k] = v
		}
	}
	for k, v := range nginxBuildOptions.Numbers {
		v.Value = flag.Int(k, v.Default, v.Desc)
		nginxBuildOptions.Numbers[k] = v
	}

	overrideUnableParseFlags()

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
	libreSSLStatic := nginxBuildOptions.Bools["libressl"].Enabled
	zlibStatic := nginxBuildOptions.Bools["zlib"].Enabled
	clear := nginxBuildOptions.Bools["clear"].Enabled
	versionPrint := nginxBuildOptions.Bools["version"].Enabled
	versionsPrint := nginxBuildOptions.Bools["versions"].Enabled
	openResty := nginxBuildOptions.Bools["openresty"].Enabled
	freenginx := nginxBuildOptions.Bools["freenginx"].Enabled
	configureOnly := nginxBuildOptions.Bools["configureonly"].Enabled
	idempotent := nginxBuildOptions.Bools["idempotent"].Enabled
	helpAll := nginxBuildOptions.Bools["help-all"].Enabled

	version := nginxBuildOptions.Values["v"].Value
	nginxConfigurePath := nginxBuildOptions.Values["c"].Value
	modulesConfPath := nginxBuildOptions.Values["m"].Value
	workParentDir := nginxBuildOptions.Values["d"].Value
	pcreVersion := nginxBuildOptions.Values["pcreversion"].Value
	openSSLVersion := nginxBuildOptions.Values["opensslversion"].Value
	libreSSLVersion := nginxBuildOptions.Values["libresslversion"].Value
	zlibVersion := nginxBuildOptions.Values["zlibversion"].Value
	openRestyVersion := nginxBuildOptions.Values["openrestyversion"].Value
	freenginxVersion := nginxBuildOptions.Values["freenginxversion"].Value
	patchOption := nginxBuildOptions.Values["patch-opt"].Value
	customSSLURL := nginxBuildOptions.Values["customssl"].Value
	customSSLName := nginxBuildOptions.Values["customsslname"].Value
	customSSLTag := nginxBuildOptions.Values["customssltag"].Value

	// Allow multiple flags for `--patch`
	{
		tmp := nginxBuildOptions.Values["patch"]
		tmp2 := multiflagPatch.String()
		tmp.Value = &tmp2
		nginxBuildOptions.Values["patch"] = tmp
	}

	// Allow multiple flags for `--add-module`
	{
		tmp := argsString["add-module"]
		tmp2 := multiflag.String()
		tmp.Value = &tmp2
		argsString["add-module"] = tmp
	}

	// Allow multiple flags for `--add-dynamic-module`
	{
		tmp := argsString["add-dynamic-module"]
		tmp2 := multiflagDynamic.String()
		tmp.Value = &tmp2
		argsString["add-dynamic-module"] = tmp
	}

	patchPath := nginxBuildOptions.Values["patch"].Value
	patchDirectives, err := parsePatchDirectives(*patchPath)
	if err != nil {
		log.Fatal(err)
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
	if *openResty && *freenginx {
		log.Fatal("select one between '-openresty' and '-freenginx'.")
	}
	// Check for conflicting SSL library options
	sslCount := 0
	if *openSSLStatic {
		sslCount++
	}
	if *libreSSLStatic {
		sslCount++
	}
	if *customSSLURL != "" {
		sslCount++
	}
	if sslCount > 1 {
		log.Fatal("select only one SSL library: '-openssl', '-libressl', or '-customssl'.")
	}
	if *openResty {
		nginxBuilder = builder.MakeBuilder(builder.ComponentOpenResty, *openRestyVersion)
	} else if *freenginx {
		nginxBuilder = builder.MakeBuilder(builder.ComponentFreenginx, *freenginxVersion)
	} else {
		nginxBuilder = builder.MakeBuilder(builder.ComponentNginx, *version)
	}
	pcreBuilder := builder.MakeLibraryBuilder(builder.ComponentPcre, *pcreVersion, *pcreStatic)
	openSSLBuilder := builder.MakeLibraryBuilder(builder.ComponentOpenSSL, *openSSLVersion, *openSSLStatic)
	libreSSLBuilder := builder.MakeLibraryBuilder(builder.ComponentLibreSSL, *libreSSLVersion, *libreSSLStatic)

	// Create custom SSL builder if URL is provided
	var customSSLBuilder builder.Builder
	if *customSSLURL != "" {
		customSSLBuilder = builder.Builder{
			Component:  builder.ComponentCustomSSL,
			Version:    "custom",
			Static:     true,
			CustomURL:  *customSSLURL,
			CustomName: *customSSLName,
			CustomTag:  *customSSLTag,
		}
	}

	zlibBuilder := builder.MakeLibraryBuilder(builder.ComponentZlib, *zlibVersion, *zlibStatic)

	if *idempotent {
		builders := []builder.Builder{
			nginxBuilder,
			pcreBuilder,
			openSSLBuilder,
			libreSSLBuilder,
			zlibBuilder,
		}
		if *customSSLURL != "" {
			builders = append(builders, customSSLBuilder)
		}

		isSame, err := builder.IsSameVersion(builders)
		if err != nil {
			log.Println("[notice]", err)
		}
		if isSame {
			log.Println("Installed nginx is same.")
			return
		}
	}

	// change default umask
	_ = syscall.Umask(0)

	versionCheck(*version)

	nginxConfigure, err := util.FileGetContents(*nginxConfigurePath)
	if err != nil {
		log.Fatal(err)
	}
	nginxConfigure, trailingComments := configure.Normalize(nginxConfigure)

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
			log.Fatalf("Failed to create working directory %s.", *workParentDir)
		}
	}

	var workDir string
	if *openResty {
		workDir = *workParentDir + "/openresty/" + *openRestyVersion
	} else if *freenginx {
		workDir = *workParentDir + "/freenginx/" + *freenginxVersion
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
			log.Fatalf("Failed to create working directory %s.", workDir)
		}
	}

	workDirAbs, err := filepath.Abs(workDir)
	if err != nil {
		log.Fatalf("Failed to resolve working directory: %v", err)
	}

	patchTargetDirs := make(map[string]string)
	addPatchTarget := func(name, relativePath string) {
		if name == "" || relativePath == "" {
			return
		}
		patchTargetDirs[strings.ToLower(name)] = filepath.Join(workDirAbs, relativePath)
	}

	addPatchTarget(defaultPatchTarget, nginxBuilder.SourcePath())
	if *openResty {
		addPatchTarget("openresty", nginxBuilder.SourcePath())
		addPatchTarget("ngx_openresty", nginxBuilder.SourcePath())
	}
	if *freenginx {
		addPatchTarget("freenginx", nginxBuilder.SourcePath())
	}
	if *pcreStatic {
		addPatchTarget("pcre", pcreBuilder.SourcePath())
		addPatchTarget("pcre2", pcreBuilder.SourcePath())
	}
	if *openSSLStatic {
		addPatchTarget("openssl", openSSLBuilder.SourcePath())
	}
	if *libreSSLStatic {
		addPatchTarget("libressl", libreSSLBuilder.SourcePath())
	}
	if *customSSLURL != "" {
		addPatchTarget("customssl", customSSLBuilder.SourcePath())
		if customSSLBuilder.CustomName != "" {
			addPatchTarget(strings.ToLower(customSSLBuilder.CustomName), customSSLBuilder.SourcePath())
		}
	}
	if *zlibStatic {
		addPatchTarget("zlib", zlibBuilder.SourcePath())
	}

	if err := validatePatchDirectives(patchDirectives, patchTargetDirs); err != nil {
		log.Fatal(err)
	}

	rootDir, err := util.SaveCurrentDir()
	if err != nil {
		log.Fatalf("Failed to get current directory: %v", err)
	}
	err = os.Chdir(workDir)
	if err != nil {
		log.Fatal(err)
	}

	if len(patchDirectives) > 0 {
		if err := resetPatchedSources(patchDirectives, patchTargetDirs); err != nil {
			log.Fatal(err)
		}
	}

	var wg sync.WaitGroup
	if *pcreStatic {
		wg.Add(1)
		go func() {
			downloadAndExtractParallel(&pcreBuilder)
			wg.Done()
		}()
	}

	if *openSSLStatic {
		wg.Add(1)
		go func() {
			downloadAndExtractParallel(&openSSLBuilder)
			wg.Done()
		}()
	}

	if *libreSSLStatic {
		wg.Add(1)
		go func() {
			downloadAndExtractParallel(&libreSSLBuilder)
			wg.Done()
		}()
	}

	if *customSSLURL != "" {
		wg.Add(1)
		go func() {
			downloadAndExtractParallel(&customSSLBuilder)
			wg.Done()
		}()
	}

	if *zlibStatic {
		wg.Add(1)
		go func() {
			downloadAndExtractParallel(&zlibBuilder)
			wg.Done()
		}()
	}

	wg.Add(1)
	go func() {
		downloadAndExtractParallel(&nginxBuilder)
		wg.Done()
	}()

	if len(modules3rd) > 0 {
		wg.Add(len(modules3rd))
		for _, m := range modules3rd {
			go func(m module3rd.Module3rd) {
				module3rd.DownloadAndExtractParallel(m)
				wg.Done()
			}(m)
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
	if err := os.Chdir(nginxBuilder.SourcePath()); err != nil {
		log.Fatalf("failed to change directory: %v", err)
	}

	var dependencies []builder.StaticLibrary
	if *pcreStatic {
		dependencies = append(dependencies, builder.MakeStaticLibrary(&pcreBuilder))
	}

	if *openSSLStatic {
		dependencies = append(dependencies, builder.MakeStaticLibrary(&openSSLBuilder))
	}

	if *libreSSLStatic {
		dependencies = append(dependencies, builder.MakeStaticLibrary(&libreSSLBuilder))
	}

	if *customSSLURL != "" {
		dependencies = append(dependencies, builder.MakeStaticLibrary(&customSSLBuilder))
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

	if *libreSSLStatic && libreSSLBuilder.IsIncludeWithOption(nginxConfigure) {
		log.Println(libreSSLBuilder.WarnMsgWithLibrary())
	}

	if *customSSLURL != "" && customSSLBuilder.IsIncludeWithOption(nginxConfigure) {
		log.Println(customSSLBuilder.WarnMsgWithLibrary())
	}

	if *zlibStatic && zlibBuilder.IsIncludeWithOption(nginxConfigure) {
		log.Println(zlibBuilder.WarnMsgWithLibrary())
	}

	configureScript := configure.Generate(nginxConfigure, modules3rd, dependencies, configureOptions, rootDir, *openResty, *jobs)
	if trailingComments != "" {
		if !strings.HasSuffix(configureScript, "\n") {
			configureScript += "\n"
		}
		configureScript += trailingComments
		if !strings.HasSuffix(configureScript, "\n") {
			configureScript += "\n"
		}
	}

	err = os.WriteFile("./nginx-configure", []byte(configureScript), 0655)
	if err != nil {
		log.Fatalf("Failed to generate configure script for %s", nginxBuilder.SourcePath())
	}

	if err := applyPatches(patchDirectives, *patchOption, rootDir, patchTargetDirs); err != nil {
		log.Fatalf("Failed to apply patch: %v", err)
	}

	// reverts source code with patch -R when the build was interrupted.
	if len(patchDirectives) > 0 {
		sigChannel := make(chan os.Signal, 1)
		signal.Notify(sigChannel, os.Interrupt)
		go func() {
			<-sigChannel
			if err := revertPatches(patchDirectives, *patchOption, rootDir, patchTargetDirs); err != nil {
				log.Printf("Failed to revert patch: %v", err)
			}
		}()
	}

	log.Printf("Configure %s.....", nginxBuilder.SourcePath())

	err = configure.Run()
	if err != nil {
		log.Printf("Failed to configure %s\n", nginxBuilder.SourcePath())
		if err := revertPatches(patchDirectives, *patchOption, rootDir, patchTargetDirs); err != nil {
			log.Printf("Failed to revert patch: %v", err)
		}
		util.PrintFatalMsg(err, "nginx-configure.log")
	}

	if *configureOnly {
		if err := revertPatches(patchDirectives, *patchOption, rootDir, patchTargetDirs); err != nil {
			log.Printf("Failed to revert patch: %v", err)
		}
		printLastMsg(workDir, nginxBuilder.SourcePath(), *openResty, *configureOnly)
		return
	}

	log.Printf("Build %s.....", nginxBuilder.SourcePath())

	if *openSSLStatic {
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
		if err := revertPatches(patchDirectives, *patchOption, rootDir, patchTargetDirs); err != nil {
			log.Printf("Failed to revert patch: %v", err)
		}
		util.PrintFatalMsg(err, "nginx-build.log")
	}

	printLastMsg(workDir, nginxBuilder.SourcePath(), *openResty, *configureOnly)
}

func parsePatchDirectives(raw string) ([]patchDirective, error) {
	if raw == "" {
		return nil, nil
	}
	entries := strings.Split(raw, ",")
	directives := make([]patchDirective, 0, len(entries))
	index := make(map[string]int)
	for _, entry := range entries {
		spec := strings.TrimSpace(entry)
		if spec == "" {
			continue
		}
		target := defaultPatchTarget
		pathSpec := spec
		if parts := strings.SplitN(spec, "=", 2); len(parts) == 2 {
			target = strings.ToLower(strings.TrimSpace(parts[0]))
			pathSpec = strings.TrimSpace(parts[1])
			if target == "" {
				return nil, fmt.Errorf("invalid patch specification %q: missing target", spec)
			}
		}
		if pathSpec == "" {
			return nil, fmt.Errorf("invalid patch specification %q: missing path", spec)
		}
		idx, ok := index[target]
		if !ok {
			directives = append(directives, patchDirective{target: target})
			idx = len(directives) - 1
			index[target] = idx
		}
		directives[idx].paths = append(directives[idx].paths, pathSpec)
	}
	return directives, nil
}

func validatePatchDirectives(directives []patchDirective, targetDirs map[string]string) error {
	for _, directive := range directives {
		if _, ok := targetDirs[directive.target]; !ok {
			available := make([]string, 0, len(targetDirs))
			for key := range targetDirs {
				available = append(available, key)
			}
			sort.Strings(available)
			return fmt.Errorf("patch target %q is not available. Available targets: %s", directive.target, strings.Join(available, ", "))
		}
	}
	return nil
}

func resetPatchedSources(directives []patchDirective, targetDirs map[string]string) error {
	if len(directives) == 0 {
		return nil
	}
	seen := make(map[string]struct{})
	for _, directive := range directives {
		targetDir := targetDirs[directive.target]
		if targetDir == "" {
			continue
		}
		if _, ok := seen[targetDir]; ok {
			continue
		}
		if util.FileExists(targetDir) {
			if err := os.RemoveAll(targetDir); err != nil {
				return fmt.Errorf("failed to remove %s: %w", targetDir, err)
			}
		}
		seen[targetDir] = struct{}{}
	}
	return nil
}

func applyPatches(directives []patchDirective, patchOption, rootDir string, targetDirs map[string]string) error {
	return runPatchCommands(directives, patchOption, rootDir, targetDirs, false)
}

func revertPatches(directives []patchDirective, patchOption, rootDir string, targetDirs map[string]string) error {
	return runPatchCommands(directives, patchOption, rootDir, targetDirs, true)
}

func runPatchCommands(directives []patchDirective, patchOption, rootDir string, targetDirs map[string]string, reverse bool) error {
	if len(directives) == 0 {
		return nil
	}
	for _, directive := range directives {
		if len(directive.paths) == 0 {
			continue
		}
		targetDir, ok := targetDirs[directive.target]
		if !ok || targetDir == "" {
			return fmt.Errorf("patch target %q is not available", directive.target)
		}
		joinedPaths := strings.Join(directive.paths, ",")
		if err := util.Patch(joinedPaths, patchOption, rootDir, targetDir, directive.target, reverse); err != nil {
			return err
		}
	}
	return nil
}
