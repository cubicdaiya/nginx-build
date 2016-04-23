package builder

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

	suite.builders[COMPONENT_NGINX] = MakeBuilder(COMPONENT_NGINX, NGINX_VERSION)
	suite.builders[COMPONENT_PCRE] = MakeBuilder(COMPONENT_PCRE, PCRE_VERSION)
	suite.builders[COMPONENT_OPENSSL] = MakeBuilder(COMPONENT_OPENSSL, OPENSSL_VERSION)
	suite.builders[COMPONENT_ZLIB] = MakeBuilder(COMPONENT_ZLIB, ZLIB_VERSION)
	suite.builders[COMPONENT_OPENRESTY] = MakeBuilder(COMPONENT_OPENRESTY, OPENRESTY_VERSION)
	suite.builders[COMPONENT_TENGINE] = MakeBuilder(COMPONENT_TENGINE, TENGINE_VERSION)
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
	assert.Equal(suite.T(), suite.builders[COMPONENT_NGINX].DownloadURL(), fmt.Sprintf("%s/nginx-%s.tar.gz", NGINX_DOWNLOAD_URL_PREFIX, NGINX_VERSION))
	assert.Equal(suite.T(), suite.builders[COMPONENT_PCRE].DownloadURL(), fmt.Sprintf("%s/pcre-%s.tar.gz", PCRE_DOWNLOAD_URL_PREFIX, PCRE_VERSION))
	assert.Equal(suite.T(), suite.builders[COMPONENT_OPENSSL].DownloadURL(), fmt.Sprintf("%s/openssl-%s.tar.gz", OPENSSL_DOWNLOAD_URL_PREFIX, OPENSSL_VERSION))
	assert.Equal(suite.T(), suite.builders[COMPONENT_ZLIB].DownloadURL(), fmt.Sprintf("%s/zlib-%s.tar.gz", ZLIB_DOWNLOAD_URL_PREFIX, ZLIB_VERSION))
	assert.Equal(suite.T(), suite.builders[COMPONENT_OPENRESTY].DownloadURL(), fmt.Sprintf("%s/openresty-%s.tar.gz", OPENRESTY_DOWNLOAD_URL_PREFIX, OPENRESTY_VERSION))
	assert.Equal(suite.T(), suite.builders[COMPONENT_TENGINE].DownloadURL(), fmt.Sprintf("%s/tengine-%s.tar.gz", TENGINE_DOWNLOAD_URL_PREFIX, TENGINE_VERSION))
}

func (suite *BuilderTestSuite) TestSourcePath() {
	assert.Equal(suite.T(), suite.builders[COMPONENT_NGINX].SourcePath(), fmt.Sprintf("nginx-%s", NGINX_VERSION))
	assert.Equal(suite.T(), suite.builders[COMPONENT_PCRE].SourcePath(), fmt.Sprintf("pcre-%s", PCRE_VERSION))
	assert.Equal(suite.T(), suite.builders[COMPONENT_OPENSSL].SourcePath(), fmt.Sprintf("openssl-%s", OPENSSL_VERSION))
	assert.Equal(suite.T(), suite.builders[COMPONENT_ZLIB].SourcePath(), fmt.Sprintf("zlib-%s", ZLIB_VERSION))
	assert.Equal(suite.T(), suite.builders[COMPONENT_OPENRESTY].SourcePath(), fmt.Sprintf("openresty-%s", OPENRESTY_VERSION))
	assert.Equal(suite.T(), suite.builders[COMPONENT_TENGINE].SourcePath(), fmt.Sprintf("tengine-%s", TENGINE_VERSION))
}

func (suite *BuilderTestSuite) TestArchivePath() {
	assert.Equal(suite.T(), suite.builders[COMPONENT_NGINX].ArchivePath(), fmt.Sprintf("nginx-%s.tar.gz", NGINX_VERSION))
	assert.Equal(suite.T(), suite.builders[COMPONENT_PCRE].ArchivePath(), fmt.Sprintf("pcre-%s.tar.gz", PCRE_VERSION))
	assert.Equal(suite.T(), suite.builders[COMPONENT_OPENSSL].ArchivePath(), fmt.Sprintf("openssl-%s.tar.gz", OPENSSL_VERSION))
	assert.Equal(suite.T(), suite.builders[COMPONENT_ZLIB].ArchivePath(), fmt.Sprintf("zlib-%s.tar.gz", ZLIB_VERSION))
	assert.Equal(suite.T(), suite.builders[COMPONENT_OPENRESTY].ArchivePath(), fmt.Sprintf("openresty-%s.tar.gz", OPENRESTY_VERSION))
	assert.Equal(suite.T(), suite.builders[COMPONENT_TENGINE].ArchivePath(), fmt.Sprintf("tengine-%s.tar.gz", TENGINE_VERSION))
}

func (suite *BuilderTestSuite) TestLogPath() {
	assert.Equal(suite.T(), suite.builders[COMPONENT_NGINX].LogPath(), fmt.Sprintf("nginx-%s.log", NGINX_VERSION))
	assert.Equal(suite.T(), suite.builders[COMPONENT_PCRE].LogPath(), fmt.Sprintf("pcre-%s.log", PCRE_VERSION))
	assert.Equal(suite.T(), suite.builders[COMPONENT_OPENSSL].LogPath(), fmt.Sprintf("openssl-%s.log", OPENSSL_VERSION))
	assert.Equal(suite.T(), suite.builders[COMPONENT_ZLIB].LogPath(), fmt.Sprintf("zlib-%s.log", ZLIB_VERSION))
	assert.Equal(suite.T(), suite.builders[COMPONENT_OPENRESTY].LogPath(), fmt.Sprintf("openresty-%s.log", OPENRESTY_VERSION))
	assert.Equal(suite.T(), suite.builders[COMPONENT_TENGINE].LogPath(), fmt.Sprintf("tengine-%s.log", TENGINE_VERSION))
}

func (suite *BuilderTestSuite) TestMakeStaticLibrary() {
	pcre := MakeStaticLibrary(&suite.builders[COMPONENT_PCRE])
	openssl := MakeStaticLibrary(&suite.builders[COMPONENT_OPENSSL])
	zlib := MakeStaticLibrary(&suite.builders[COMPONENT_ZLIB])

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

func TestBuilderTestSuite(t *testing.T) {
	suite.Run(t, new(BuilderTestSuite))
}
