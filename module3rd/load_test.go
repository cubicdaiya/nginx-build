package module3rd

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
	modules3rdConf := "../config/modules.cfg.example"
	modules3rd, err := Load(modules3rdConf)
	if err != nil {
		log.Fatalf("Failed to load %s\n", modules3rdConf)
	}
	suite.modules3rd = modules3rd
}

func (suite *Module3rdTestSuite) TestModules3rd() {
	for _, m := range suite.modules3rd {
		switch m.Name {
		//		case "echo-nginx-module":
		//			assert.Equal(suite.T(), m.Name, "echo-nginx-module")
		//			assert.Equal(suite.T(), m.Form, "git")
		//			assert.Equal(suite.T(), m.Url, "https://github.com/openresty/echo-nginx-module.git")
		//			assert.Equal(suite.T(), m.Rev, "v0.58")
		//			assert.Equal(suite.T(), m.Dynamic, false)
		case "headers-more-nginx-module":
			assert.Equal(suite.T(), m.Name, "headers-more-nginx-module")
			assert.Equal(suite.T(), m.Form, "git")
			assert.Equal(suite.T(), m.Url, "https://github.com/openresty/headers-more-nginx-module.git")
			assert.Equal(suite.T(), m.Rev, "v0.32")
			assert.Equal(suite.T(), m.Dynamic, false)
		case "ngx_devel_kit":
			assert.Equal(suite.T(), m.Name, "ngx_devel_kit")
			assert.Equal(suite.T(), m.Form, "git")
			assert.Equal(suite.T(), m.Url, "https://github.com/simpl/ngx_devel_kit")
			assert.Equal(suite.T(), m.Rev, "v0.3.0")
			assert.Equal(suite.T(), m.Dynamic, false)
		case "ngx_small_light":
			assert.Equal(suite.T(), m.Name, "ngx_small_light")
			assert.Equal(suite.T(), m.Form, "git")
			assert.Equal(suite.T(), m.Url, "https://github.com/cubicdaiya/ngx_small_light")
			assert.Equal(suite.T(), m.Rev, "v0.9.2")
			assert.Equal(suite.T(), m.Shprov, "./setup --with-gd")
			assert.Equal(suite.T(), m.Dynamic, true)
		}
	}
}

func TestModule3rdTestSuite(t *testing.T) {
	suite.Run(t, new(Module3rdTestSuite))
}
