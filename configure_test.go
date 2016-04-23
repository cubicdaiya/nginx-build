package main

import (
	"fmt"
	"log"
	"runtime"
	"testing"

	"github.com/cubicdaiya/nginx-build/builder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ConfiguregenTestSuite struct {
	suite.Suite
	builders   []builder.Builder
	modules3rd []Module3rd
}

func (suite *ConfiguregenTestSuite) SetupTest() {
	suite.builders = make([]builder.Builder, builder.COMPONENT_MAX)

	suite.builders[builder.COMPONENT_PCRE] = builder.MakeBuilder(builder.COMPONENT_PCRE, builder.PCRE_VERSION)
	suite.builders[builder.COMPONENT_OPENSSL] = builder.MakeBuilder(builder.COMPONENT_OPENSSL, builder.OPENSSL_VERSION)
	suite.builders[builder.COMPONENT_ZLIB] = builder.MakeBuilder(builder.COMPONENT_ZLIB, builder.ZLIB_VERSION)

	modules3rdConf := "./config/modules.cfg.example"
	modules3rd, err := loadModules3rdFile(modules3rdConf)
	if err != nil {
		log.Fatalf("Failed to load %s\n", modules3rdConf)
	}
	suite.modules3rd = modules3rd
}

func (suite *ConfiguregenTestSuite) TestConfiguregenModule3rd() {
	configure_modules3rd := configureGenModule3rd(suite.modules3rd)

	assert.Contains(suite.T(), configure_modules3rd, "-add-module=../echo-nginx-module")
	assert.Contains(suite.T(), configure_modules3rd, "-add-module=../headers-more-nginx-module")
	assert.Contains(suite.T(), configure_modules3rd, "-add-module=../ngx_devel_kit")
	assert.Contains(suite.T(), configure_modules3rd, "-add-dynamic-module=../ngx_small_light")
}

func (suite *ConfiguregenTestSuite) TestConfiguregenDefault() {
	var configureOptions ConfigureOptions
	configureScript := configureGen("", []Module3rd{}, []builder.StaticLibrary{}, configureOptions, "")

	if runtime.GOOS == "darwin" {
		assert.Contains(suite.T(), configureScript, "--with-cc-opt=\"-Wno-deprecated-declarations\" \\")
	}
}

func (suite *ConfiguregenTestSuite) TestConfiguregenWithStaticLibraries() {
	var dependencies []builder.StaticLibrary
	dependencies = append(dependencies, builder.MakeStaticLibrary(&suite.builders[builder.COMPONENT_PCRE]))
	dependencies = append(dependencies, builder.MakeStaticLibrary(&suite.builders[builder.COMPONENT_OPENSSL]))
	dependencies = append(dependencies, builder.MakeStaticLibrary(&suite.builders[builder.COMPONENT_ZLIB]))
	var configureOptions ConfigureOptions
	configureScript := configureGen("", []Module3rd{}, dependencies, configureOptions, "")

	assert.Contains(suite.T(), configureScript, "--with-http_ssl_module")
	assert.Contains(suite.T(), configureScript, fmt.Sprintf("--with-pcre=../pcre-%s \\\n", builder.PCRE_VERSION))
	assert.Contains(suite.T(), configureScript, fmt.Sprintf("--with-openssl=../openssl-%s \\\n", builder.OPENSSL_VERSION))
	assert.Contains(suite.T(), configureScript, fmt.Sprintf("--with-zlib=../zlib-%s \\\n", builder.ZLIB_VERSION))
}

func TestConfiguregenTestSuite(t *testing.T) {
	suite.Run(t, new(ConfiguregenTestSuite))
}
