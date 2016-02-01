package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type BuilderTestSuite struct {
	suite.Suite
	builders []Builder
}

func (suite *BuilderTestSuite) SetupTest() {
	suite.builders = make([]Builder, COMPONENT_MAX)

	suite.builders[COMPONENT_NGINX] = makeBuilder(COMPONENT_NGINX, NGINX_VERSION)
	suite.builders[COMPONENT_PCRE] = makeBuilder(COMPONENT_PCRE, PCRE_VERSION)
	suite.builders[COMPONENT_OPENSSL] = makeBuilder(COMPONENT_OPENSSL, OPENSSL_VERSION)
	suite.builders[COMPONENT_ZLIB] = makeBuilder(COMPONENT_ZLIB, ZLIB_VERSION)
	suite.builders[COMPONENT_OPENRESTY] = makeBuilder(COMPONENT_OPENRESTY, OPENRESTY_VERSION)
	suite.builders[COMPONENT_TENGINE] = makeBuilder(COMPONENT_TENGINE, TENGINE_VERSION)
}

func (suite *BuilderTestSuite) TestName() {
	assert.Equal(suite.T(), suite.builders[COMPONENT_NGINX].name(), "nginx")
	assert.Equal(suite.T(), suite.builders[COMPONENT_PCRE].name(), "pcre")
	assert.Equal(suite.T(), suite.builders[COMPONENT_OPENSSL].name(), "openssl")
	assert.Equal(suite.T(), suite.builders[COMPONENT_ZLIB].name(), "zlib")
	assert.Equal(suite.T(), suite.builders[COMPONENT_OPENRESTY].name(), "openresty")
	assert.Equal(suite.T(), suite.builders[COMPONENT_TENGINE].name(), "tengine")
}

func (suite *BuilderTestSuite) TestOption() {
	assert.Equal(suite.T(), suite.builders[COMPONENT_PCRE].option(), "--with-pcre")
	assert.Equal(suite.T(), suite.builders[COMPONENT_OPENSSL].option(), "--with-openssl")
	assert.Equal(suite.T(), suite.builders[COMPONENT_ZLIB].option(), "--with-zlib")
}

func (suite *BuilderTestSuite) TestDownloadURL() {
	assert.Equal(suite.T(), suite.builders[COMPONENT_NGINX].downloadURL(), fmt.Sprintf("%s/nginx-%s.tar.gz", NGINX_DOWNLOAD_URL_PREFIX, NGINX_VERSION))
	assert.Equal(suite.T(), suite.builders[COMPONENT_PCRE].downloadURL(), fmt.Sprintf("%s/pcre-%s.tar.gz", PCRE_DOWNLOAD_URL_PREFIX, PCRE_VERSION))
	assert.Equal(suite.T(), suite.builders[COMPONENT_OPENSSL].downloadURL(), fmt.Sprintf("%s/openssl-%s.tar.gz", OPENSSL_DOWNLOAD_URL_PREFIX, OPENSSL_VERSION))
	assert.Equal(suite.T(), suite.builders[COMPONENT_ZLIB].downloadURL(), fmt.Sprintf("%s/zlib-%s.tar.gz", ZLIB_DOWNLOAD_URL_PREFIX, ZLIB_VERSION))
	assert.Equal(suite.T(), suite.builders[COMPONENT_OPENRESTY].downloadURL(), fmt.Sprintf("%s/openresty-%s.tar.gz", OPENRESTY_DOWNLOAD_URL_PREFIX, OPENRESTY_VERSION))
	assert.Equal(suite.T(), suite.builders[COMPONENT_TENGINE].downloadURL(), fmt.Sprintf("%s/tengine-%s.tar.gz", TENGINE_DOWNLOAD_URL_PREFIX, TENGINE_VERSION))
}

func (suite *BuilderTestSuite) TestSourcePath() {
	assert.Equal(suite.T(), suite.builders[COMPONENT_NGINX].sourcePath(), fmt.Sprintf("nginx-%s", NGINX_VERSION))
	assert.Equal(suite.T(), suite.builders[COMPONENT_PCRE].sourcePath(), fmt.Sprintf("pcre-%s", PCRE_VERSION))
	assert.Equal(suite.T(), suite.builders[COMPONENT_OPENSSL].sourcePath(), fmt.Sprintf("openssl-%s", OPENSSL_VERSION))
	assert.Equal(suite.T(), suite.builders[COMPONENT_ZLIB].sourcePath(), fmt.Sprintf("zlib-%s", ZLIB_VERSION))
	assert.Equal(suite.T(), suite.builders[COMPONENT_OPENRESTY].sourcePath(), fmt.Sprintf("openresty-%s", OPENRESTY_VERSION))
	assert.Equal(suite.T(), suite.builders[COMPONENT_TENGINE].sourcePath(), fmt.Sprintf("tengine-%s", TENGINE_VERSION))
}

func (suite *BuilderTestSuite) TestArchivePath() {
	assert.Equal(suite.T(), suite.builders[COMPONENT_NGINX].archivePath(), fmt.Sprintf("nginx-%s.tar.gz", NGINX_VERSION))
	assert.Equal(suite.T(), suite.builders[COMPONENT_PCRE].archivePath(), fmt.Sprintf("pcre-%s.tar.gz", PCRE_VERSION))
	assert.Equal(suite.T(), suite.builders[COMPONENT_OPENSSL].archivePath(), fmt.Sprintf("openssl-%s.tar.gz", OPENSSL_VERSION))
	assert.Equal(suite.T(), suite.builders[COMPONENT_ZLIB].archivePath(), fmt.Sprintf("zlib-%s.tar.gz", ZLIB_VERSION))
	assert.Equal(suite.T(), suite.builders[COMPONENT_OPENRESTY].archivePath(), fmt.Sprintf("openresty-%s.tar.gz", OPENRESTY_VERSION))
	assert.Equal(suite.T(), suite.builders[COMPONENT_TENGINE].archivePath(), fmt.Sprintf("tengine-%s.tar.gz", TENGINE_VERSION))
}

func (suite *BuilderTestSuite) TestLogPath() {
	assert.Equal(suite.T(), suite.builders[COMPONENT_NGINX].logPath(), fmt.Sprintf("nginx-%s.log", NGINX_VERSION))
	assert.Equal(suite.T(), suite.builders[COMPONENT_PCRE].logPath(), fmt.Sprintf("pcre-%s.log", PCRE_VERSION))
	assert.Equal(suite.T(), suite.builders[COMPONENT_OPENSSL].logPath(), fmt.Sprintf("openssl-%s.log", OPENSSL_VERSION))
	assert.Equal(suite.T(), suite.builders[COMPONENT_ZLIB].logPath(), fmt.Sprintf("zlib-%s.log", ZLIB_VERSION))
	assert.Equal(suite.T(), suite.builders[COMPONENT_OPENRESTY].logPath(), fmt.Sprintf("openresty-%s.log", OPENRESTY_VERSION))
	assert.Equal(suite.T(), suite.builders[COMPONENT_TENGINE].logPath(), fmt.Sprintf("tengine-%s.log", TENGINE_VERSION))
}

func TestBuilderTestSuite(t *testing.T) {
	suite.Run(t, new(BuilderTestSuite))
}
