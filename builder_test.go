package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
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
}

func (suite *BuilderTestSuite) TestName() {
	assert.Equal(suite.T(), suite.builders[COMPONENT_NGINX].name(), "nginx")
	assert.Equal(suite.T(), suite.builders[COMPONENT_PCRE].name(), "pcre")
	assert.Equal(suite.T(), suite.builders[COMPONENT_OPENSSL].name(), "openssl")
	assert.Equal(suite.T(), suite.builders[COMPONENT_ZLIB].name(), "zlib")
}

func (suite *BuilderTestSuite) TestOption() {
	assert.Equal(suite.T(), suite.builders[COMPONENT_NGINX].option(), "--with-nginx")
	assert.Equal(suite.T(), suite.builders[COMPONENT_PCRE].option(), "--with-pcre")
	assert.Equal(suite.T(), suite.builders[COMPONENT_OPENSSL].option(), "--with-openssl")
	assert.Equal(suite.T(), suite.builders[COMPONENT_ZLIB].option(), "--with-zlib")
}

func (suite *BuilderTestSuite) TestDownloadURL() {
	assert.Equal(suite.T(), suite.builders[COMPONENT_NGINX].downloadURL(), fmt.Sprintf("%s/nginx-%s.tar.gz", NGINX_DOWNLOAD_URL_PREFIX, NGINX_VERSION))
	assert.Equal(suite.T(), suite.builders[COMPONENT_PCRE].downloadURL(), fmt.Sprintf("%s/pcre-%s.tar.gz", PCRE_DOWNLOAD_URL_PREFIX, PCRE_VERSION))
	assert.Equal(suite.T(), suite.builders[COMPONENT_OPENSSL].downloadURL(), fmt.Sprintf("%s/openssl-%s.tar.gz", OPENSSL_DOWNLOAD_URL_PREFIX, OPENSSL_VERSION))
}

func (suite *BuilderTestSuite) TestSourcePath() {
	assert.Equal(suite.T(), suite.builders[COMPONENT_NGINX].sourcePath(), fmt.Sprintf("nginx-%s", NGINX_VERSION))
	assert.Equal(suite.T(), suite.builders[COMPONENT_PCRE].sourcePath(), fmt.Sprintf("pcre-%s", PCRE_VERSION))
	assert.Equal(suite.T(), suite.builders[COMPONENT_OPENSSL].sourcePath(), fmt.Sprintf("openssl-%s", OPENSSL_VERSION))
	assert.Equal(suite.T(), suite.builders[COMPONENT_ZLIB].sourcePath(), fmt.Sprintf("zlib-%s", ZLIB_VERSION))
}

func (suite *BuilderTestSuite) TestArchivePath() {
	assert.Equal(suite.T(), suite.builders[COMPONENT_NGINX].archivePath(), fmt.Sprintf("nginx-%s.tar.gz", NGINX_VERSION))
	assert.Equal(suite.T(), suite.builders[COMPONENT_PCRE].archivePath(), fmt.Sprintf("pcre-%s.tar.gz", PCRE_VERSION))
	assert.Equal(suite.T(), suite.builders[COMPONENT_OPENSSL].archivePath(), fmt.Sprintf("openssl-%s.tar.gz", OPENSSL_VERSION))
	assert.Equal(suite.T(), suite.builders[COMPONENT_ZLIB].archivePath(), fmt.Sprintf("zlib-%s.tar.gz", ZLIB_VERSION))
}

func TestBuilderTestSuite(t *testing.T) {
	suite.Run(t, new(BuilderTestSuite))
}
