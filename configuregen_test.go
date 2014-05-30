package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
	"runtime"
	"testing"
)

type ConfiguregenTestSuite struct {
	suite.Suite
	builders   []Builder
	modules3rd []Module3rd
}

func (suite *ConfiguregenTestSuite) SetupTest() {
	suite.builders = make([]Builder, COMPONENT_MAX)

	suite.builders[COMPONENT_PCRE] = makeBuilder(COMPONENT_PCRE, PCRE_VERSION)
	suite.builders[COMPONENT_OPENSSL] = makeBuilder(COMPONENT_OPENSSL, OPENSSL_VERSION)
	suite.builders[COMPONENT_ZLIB] = makeBuilder(COMPONENT_ZLIB, ZLIB_VERSION)

	modules3rdConf := "./config/modules.cfg.example"
	modules3rd, err := loadModules3rdFile(modules3rdConf)
	if err != nil {
		fmt.Printf("Failed to load %s\n", modules3rdConf)
		os.Exit(1)
	}
	suite.modules3rd = modules3rd
}

func (suite *ConfiguregenTestSuite) TestMakeStaticLibrary() {
	pcre := makeStaticLibrary(&suite.builders[COMPONENT_PCRE])
	openssl := makeStaticLibrary(&suite.builders[COMPONENT_OPENSSL])
	zlib := makeStaticLibrary(&suite.builders[COMPONENT_ZLIB])

	assert.Equal(suite.T(), pcre.Name, suite.builders[COMPONENT_PCRE].name())
	assert.Equal(suite.T(), pcre.Version, PCRE_VERSION)
	assert.Equal(suite.T(), pcre.Option, suite.builders[COMPONENT_PCRE].option())

	assert.Equal(suite.T(), openssl.Name, suite.builders[COMPONENT_OPENSSL].name())
	assert.Equal(suite.T(), openssl.Version, OPENSSL_VERSION)
	assert.Equal(suite.T(), openssl.Option, suite.builders[COMPONENT_OPENSSL].option())

	assert.Equal(suite.T(), zlib.Name, suite.builders[COMPONENT_ZLIB].name())
	assert.Equal(suite.T(), zlib.Version, ZLIB_VERSION)
	assert.Equal(suite.T(), zlib.Option, suite.builders[COMPONENT_ZLIB].option())
}

func (suite *ConfiguregenTestSuite) TestConfiguregenModule3rd() {
	configure_modules3rd := configureGenModule3rd(suite.modules3rd)

	assert.Contains(suite.T(), configure_modules3rd, "-add-module=../echo-nginx-module")
	assert.Contains(suite.T(), configure_modules3rd, "-add-module=../headers-more-nginx-module")
	assert.Contains(suite.T(), configure_modules3rd, "-add-module=../ngx_devel_kit")
	assert.Contains(suite.T(), configure_modules3rd, "-add-module=../ngx_info")
	assert.Contains(suite.T(), configure_modules3rd, "-add-module=../ngx_small_light")
}

func (suite *ConfiguregenTestSuite) TestConfiguregenDefault() {
	configureScript := configureGen("", []Module3rd{}, []StaticLibrary{})

	if runtime.GOOS == "darwin" {
		assert.Contains(suite.T(), configureScript, "--with-cc-opt=\"-Wno-deprecated-declarations\" \\")
	}
}

func (suite *ConfiguregenTestSuite) TestConfiguregenWithStaticLibraries() {
	var dependencies []StaticLibrary
	dependencies = append(dependencies, makeStaticLibrary(&suite.builders[COMPONENT_PCRE]))
	dependencies = append(dependencies, makeStaticLibrary(&suite.builders[COMPONENT_OPENSSL]))
	dependencies = append(dependencies, makeStaticLibrary(&suite.builders[COMPONENT_ZLIB]))
	configureScript := configureGen("", []Module3rd{}, dependencies)

	assert.Contains(suite.T(), configureScript, fmt.Sprintf("--with-pcre=../pcre-%s \\\n", PCRE_VERSION))
	assert.Contains(suite.T(), configureScript, fmt.Sprintf("--with-openssl=../openssl-%s \\\n", OPENSSL_VERSION))
	assert.Contains(suite.T(), configureScript, fmt.Sprintf("--with-zlib=../zlib-%s \\\n", ZLIB_VERSION))
}

func TestConfiguregenTestSuite(t *testing.T) {
	suite.Run(t, new(ConfiguregenTestSuite))
}
