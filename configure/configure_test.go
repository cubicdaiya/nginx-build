package configure

import (
	"fmt"
	"strings"
	"testing"

	"github.com/cubicdaiya/nginx-build/builder"
	"github.com/cubicdaiya/nginx-build/module3rd"
)

func setupBuilders(t *testing.T) []builder.Builder {
	builders := make([]builder.Builder, builder.ComponentMax)
	builders[builder.ComponentPcre] = builder.MakeBuilder(builder.ComponentPcre, builder.PcreVersion)
	builders[builder.ComponentOpenSSL] = builder.MakeBuilder(builder.ComponentOpenSSL, builder.OpenSSLVersion)
	builders[builder.ComponentZlib] = builder.MakeBuilder(builder.ComponentZlib, builder.ZlibVersion)
	return builders
}

func setupModules3rd(t *testing.T) []module3rd.Module3rd {
	modules3rdConf := "../config/modules.json.example"
	modules3rd, err := module3rd.Load(modules3rdConf)
	if err != nil {
		t.Fatalf("Failed to load %s\n", modules3rdConf)
	}
	return modules3rd
}

func TestConfiguregenModule3rd(t *testing.T) {
	modules3rd := setupModules3rd(t)

	configureModules3rd := generateForModule3rd(modules3rd)

	wantedOptions := []string{
		"-add-module=../ngx_http_hello_world",
	}

	for _, want := range wantedOptions {
		if !strings.Contains(configureModules3rd, want) {
			t.Fatalf("configure string does not contain wanted option: %v", want)
		}
	}
}

func TestConfiguregenWithStaticLibraries(t *testing.T) {

	builders := setupBuilders(t)

	var dependencies []builder.StaticLibrary
	dependencies = append(dependencies, builder.MakeStaticLibrary(&builders[builder.ComponentPcre]))
	dependencies = append(dependencies, builder.MakeStaticLibrary(&builders[builder.ComponentOpenSSL]))
	dependencies = append(dependencies, builder.MakeStaticLibrary(&builders[builder.ComponentZlib]))
	var configureOptions Options
	configureScript := Generate("", []module3rd.Module3rd{}, dependencies, configureOptions, "", false, 1)

	wantedOptions := []string{
		"--with-http_ssl_module",
		fmt.Sprintf("--with-pcre=../pcre2-%s \\\n", builder.PcreVersion),
		fmt.Sprintf("--with-openssl=../openssl-%s \\\n", builder.OpenSSLVersion),
		fmt.Sprintf("--with-zlib=../zlib-%s \\\n", builder.ZlibVersion),
	}

	for _, want := range wantedOptions {
		if !strings.Contains(configureScript, want) {
			t.Fatalf("configure script does not contain wanted option: %v", want)
		}

	}
}

func TestGeneratePreservesTrailingComments(t *testing.T) {
	script := "#!/bin/sh\n\n./configure \\\n--prefix=/usr/local \\\n# --with-http_ssl_module\n"
	base, comments := Normalize(script)
	if base == "" {
		t.Fatalf("expected base configure content, got empty string")
	}
	if comments == "" {
		t.Fatalf("expected trailing comments to be captured")
	}
	enabled := true
	options := Options{
		Values: MakeArgsString(),
		Bools:  MakeArgsBool(),
	}
	options.Bools["stub"] = OptionBool{Name: "--with-http_stub_status_module", Enabled: &enabled}
	generated := Generate(base, nil, nil, options, "", false, 1)
	if comments != "" {
		if !strings.HasSuffix(generated, "\n") {
			generated += "\n"
		}
		generated += comments
		if !strings.HasSuffix(generated, "\n") {
			generated += "\n"
		}
	}
	if !strings.Contains(generated, "--with-http_stub_status_module") {
		t.Fatalf("generated script missing appended option: %s", generated)
	}
	if strings.Contains(generated, "# --with-http_stub_status_module") {
		t.Fatalf("appended option is still commented out: %s", generated)
	}
	if !strings.Contains(generated, "# --with-http_ssl_module") {
		t.Fatalf("trailing comment was not preserved: %s", generated)
	}
}
