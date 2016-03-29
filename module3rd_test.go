package main

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type Module3rdTestSuite struct {
	suite.Suite
	modules3rd []Module3rd
}

func (suite *Module3rdTestSuite) SetupTest() {
	modules3rdConf := "./config/modules.cfg.example"
	modules3rd, err := loadModules3rdFile(modules3rdConf)
	if err != nil {
		log.Fatalf("Failed to load %s\n", modules3rdConf)
	}
	suite.modules3rd = modules3rd
}

func (suite *Module3rdTestSuite) TestModules3rd() {
	for _, m := range suite.modules3rd {
		switch m.Name {
		case "echo-nginx-module":
			assert.Equal(suite.T(), m.Name, "echo-nginx-module")
			assert.Equal(suite.T(), m.Url, "https://github.com/openresty/echo-nginx-module.git")
			assert.Equal(suite.T(), m.Rev, "v0.58")
		case "headers-more-nginx-module":
			assert.Equal(suite.T(), m.Name, "headers-more-nginx-module")
			assert.Equal(suite.T(), m.Url, "https://github.com/openresty/headers-more-nginx-module.git")
			assert.Equal(suite.T(), m.Rev, "v0.26")
		case "ngx_devel_kit":
			assert.Equal(suite.T(), m.Name, "ngx_devel_kit")
			assert.Equal(suite.T(), m.Url, "https://github.com/simpl/ngx_devel_kit")
			assert.Equal(suite.T(), m.Rev, "v0.2.19")
		case "ngx_small_light":
			assert.Equal(suite.T(), m.Name, "ngx_small_light")
			assert.Equal(suite.T(), m.Url, "https://github.com/cubicdaiya/ngx_small_light")
			assert.Equal(suite.T(), m.Rev, "v0.6.15")
			assert.Equal(suite.T(), m.Shprov, "./setup --with-gd")
		}
	}
}

func TestModule3rdTestSuite(t *testing.T) {
	suite.Run(t, new(Module3rdTestSuite))
}
