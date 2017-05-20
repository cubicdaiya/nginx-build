package configure

import (
	"fmt"
	"log"
	"testing"

	"github.com/cubicdaiya/nginx-build/builder"
	"github.com/cubicdaiya/nginx-build/module3rd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ConfiguregenTestSuite struct {
	suite.Suite
	builders   []builder.Builder
	modules3rd []module3rd.Module3rd
}

func (suite *ConfiguregenTestSuite) SetupTest() {
	suite.builders = make([]builder.Builder, builder.ComponentMax)

	suite.builders[builder.ComponentPcre] = builder.MakeBuilder(builder.ComponentPcre, builder.PcreVersion)
	suite.builders[builder.ComponentOpenSSL] = builder.MakeBuilder(builder.ComponentOpenSSL, builder.OpenSSLVersion)
	suite.builders[builder.ComponentZlib] = builder.MakeBuilder(builder.ComponentZlib, builder.ZlibVersion)

	modules3rdConf := "../config/modules.cfg.example"
	modules3rd, err := module3rd.Load(modules3rdConf)
	if err != nil {
		log.Fatalf("Failed to load %s\n", modules3rdConf)
	}
	suite.modules3rd = modules3rd
}

func (suite *ConfiguregenTestSuite) TestConfiguregenModule3rd() {
	configure_modules3rd := generateForModule3rd(suite.modules3rd)

	//assert.Contains(suite.T(), configure_modules3rd, "-add-module=../echo-nginx-module")
	assert.Contains(suite.T(), configure_modules3rd, "-add-module=../headers-more-nginx-module")
	assert.Contains(suite.T(), configure_modules3rd, "-add-module=../ngx_devel_kit")
	assert.Contains(suite.T(), configure_modules3rd, "-add-dynamic-module=../ngx_small_light")
}

func (suite *ConfiguregenTestSuite) TestConfiguregenWithStaticLibraries() {
	var dependencies []builder.StaticLibrary
	dependencies = append(dependencies, builder.MakeStaticLibrary(&suite.builders[builder.ComponentPcre]))
	dependencies = append(dependencies, builder.MakeStaticLibrary(&suite.builders[builder.ComponentOpenSSL]))
	dependencies = append(dependencies, builder.MakeStaticLibrary(&suite.builders[builder.ComponentZlib]))
	var configureOptions Options
	configureScript := Generate("", []module3rd.Module3rd{}, dependencies, configureOptions, "", false, 1)

	assert.Contains(suite.T(), configureScript, "--with-http_ssl_module")
	assert.Contains(suite.T(), configureScript, fmt.Sprintf("--with-pcre=../pcre-%s \\\n", builder.PcreVersion))
	assert.Contains(suite.T(), configureScript, fmt.Sprintf("--with-openssl=../openssl-%s \\\n", builder.OpenSSLVersion))
	assert.Contains(suite.T(), configureScript, fmt.Sprintf("--with-zlib=../zlib-%s \\\n", builder.ZlibVersion))
}

func TestConfiguregenTestSuite(t *testing.T) {
	suite.Run(t, new(ConfiguregenTestSuite))
}
