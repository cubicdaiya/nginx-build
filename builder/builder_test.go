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
	suite.builders = make([]Builder, ComponentMax)

	suite.builders[ComponentNginx] = MakeBuilder(ComponentNginx, NginxVersion)
	suite.builders[ComponentPcre] = MakeLibraryBuilder(ComponentPcre, PcreVersion, false)
	suite.builders[ComponentOpenSSL] = MakeLibraryBuilder(ComponentOpenSSL, OpenSSLVersion, true)
	suite.builders[ComponentZlib] = MakeLibraryBuilder(ComponentZlib, ZlibVersion, false)
	suite.builders[ComponentOpenResty] = MakeBuilder(ComponentOpenResty, OpenRestyVersion)
	suite.builders[ComponentTengine] = MakeBuilder(ComponentTengine, TengineVersion)
}

func (suite *BuilderTestSuite) TestName() {
	assert.Equal(suite.T(), suite.builders[ComponentNginx].name(), "nginx")
	assert.Equal(suite.T(), suite.builders[ComponentPcre].name(), "pcre")
	assert.Equal(suite.T(), suite.builders[ComponentOpenSSL].name(), "openssl")
	assert.Equal(suite.T(), suite.builders[ComponentZlib].name(), "zlib")
	assert.Equal(suite.T(), suite.builders[ComponentOpenResty].name(), "openresty")
	assert.Equal(suite.T(), suite.builders[ComponentTengine].name(), "tengine")
}

func (suite *BuilderTestSuite) TestOption() {
	assert.Equal(suite.T(), suite.builders[ComponentPcre].option(), "--with-pcre")
	assert.Equal(suite.T(), suite.builders[ComponentOpenSSL].option(), "--with-openssl")
	assert.Equal(suite.T(), suite.builders[ComponentZlib].option(), "--with-zlib")
}

func (suite *BuilderTestSuite) TestDownloadURL() {
	assert.Equal(suite.T(), suite.builders[ComponentNginx].DownloadURL(), fmt.Sprintf("%s/nginx-%s.tar.gz", NginxDownloadURLPrefix, NginxVersion))
	assert.Equal(suite.T(), suite.builders[ComponentPcre].DownloadURL(), fmt.Sprintf("%s/pcre-%s.tar.gz", PcreDownloadURLPrefix, PcreVersion))
	assert.Equal(suite.T(), suite.builders[ComponentOpenSSL].DownloadURL(), fmt.Sprintf("%s/openssl-%s.tar.gz", OpenSSLDownloadURLPrefix, OpenSSLVersion))
	assert.Equal(suite.T(), suite.builders[ComponentZlib].DownloadURL(), fmt.Sprintf("%s/zlib-%s.tar.gz", ZlibDownloadURLPrefix, ZlibVersion))
	assert.Equal(suite.T(), suite.builders[ComponentOpenResty].DownloadURL(), fmt.Sprintf("%s/openresty-%s.tar.gz", OpenRestyDownloadURLPrefix, OpenRestyVersion))
	assert.Equal(suite.T(), suite.builders[ComponentTengine].DownloadURL(), fmt.Sprintf("%s/tengine-%s.tar.gz", TengineDownloadURLPrefix, TengineVersion))
}

func (suite *BuilderTestSuite) TestSourcePath() {
	assert.Equal(suite.T(), suite.builders[ComponentNginx].SourcePath(), fmt.Sprintf("nginx-%s", NginxVersion))
	assert.Equal(suite.T(), suite.builders[ComponentPcre].SourcePath(), fmt.Sprintf("pcre-%s", PcreVersion))
	assert.Equal(suite.T(), suite.builders[ComponentOpenSSL].SourcePath(), fmt.Sprintf("openssl-%s", OpenSSLVersion))
	assert.Equal(suite.T(), suite.builders[ComponentZlib].SourcePath(), fmt.Sprintf("zlib-%s", ZlibVersion))
	assert.Equal(suite.T(), suite.builders[ComponentOpenResty].SourcePath(), fmt.Sprintf("openresty-%s", OpenRestyVersion))
	assert.Equal(suite.T(), suite.builders[ComponentTengine].SourcePath(), fmt.Sprintf("tengine-%s", TengineVersion))
}

func (suite *BuilderTestSuite) TestArchivePath() {
	assert.Equal(suite.T(), suite.builders[ComponentNginx].ArchivePath(), fmt.Sprintf("nginx-%s.tar.gz", NginxVersion))
	assert.Equal(suite.T(), suite.builders[ComponentPcre].ArchivePath(), fmt.Sprintf("pcre-%s.tar.gz", PcreVersion))
	assert.Equal(suite.T(), suite.builders[ComponentOpenSSL].ArchivePath(), fmt.Sprintf("openssl-%s.tar.gz", OpenSSLVersion))
	assert.Equal(suite.T(), suite.builders[ComponentZlib].ArchivePath(), fmt.Sprintf("zlib-%s.tar.gz", ZlibVersion))
	assert.Equal(suite.T(), suite.builders[ComponentOpenResty].ArchivePath(), fmt.Sprintf("openresty-%s.tar.gz", OpenRestyVersion))
	assert.Equal(suite.T(), suite.builders[ComponentTengine].ArchivePath(), fmt.Sprintf("tengine-%s.tar.gz", TengineVersion))
}

func (suite *BuilderTestSuite) TestLogPath() {
	assert.Equal(suite.T(), suite.builders[ComponentNginx].LogPath(), fmt.Sprintf("nginx-%s.log", NginxVersion))
	assert.Equal(suite.T(), suite.builders[ComponentPcre].LogPath(), fmt.Sprintf("pcre-%s.log", PcreVersion))
	assert.Equal(suite.T(), suite.builders[ComponentOpenSSL].LogPath(), fmt.Sprintf("openssl-%s.log", OpenSSLVersion))
	assert.Equal(suite.T(), suite.builders[ComponentZlib].LogPath(), fmt.Sprintf("zlib-%s.log", ZlibVersion))
	assert.Equal(suite.T(), suite.builders[ComponentOpenResty].LogPath(), fmt.Sprintf("openresty-%s.log", OpenRestyVersion))
	assert.Equal(suite.T(), suite.builders[ComponentTengine].LogPath(), fmt.Sprintf("tengine-%s.log", TengineVersion))
}

func (suite *BuilderTestSuite) TestLibrary() {
	assert.Equal(suite.T(), suite.builders[ComponentPcre].Static, false)
	assert.Equal(suite.T(), suite.builders[ComponentOpenSSL].Static, true)
	assert.Equal(suite.T(), suite.builders[ComponentZlib].Static, false)
}

func (suite *BuilderTestSuite) TestMakeStaticLibrary() {
	pcre := MakeStaticLibrary(&suite.builders[ComponentPcre])
	openssl := MakeStaticLibrary(&suite.builders[ComponentOpenSSL])
	zlib := MakeStaticLibrary(&suite.builders[ComponentZlib])

	assert.Equal(suite.T(), pcre.Name, suite.builders[ComponentPcre].name())
	assert.Equal(suite.T(), pcre.Version, PcreVersion)
	assert.Equal(suite.T(), pcre.Option, suite.builders[ComponentPcre].option())

	assert.Equal(suite.T(), openssl.Name, suite.builders[ComponentOpenSSL].name())
	assert.Equal(suite.T(), openssl.Version, OpenSSLVersion)
	assert.Equal(suite.T(), openssl.Option, suite.builders[ComponentOpenSSL].option())

	assert.Equal(suite.T(), zlib.Name, suite.builders[ComponentZlib].name())
	assert.Equal(suite.T(), zlib.Version, ZlibVersion)
	assert.Equal(suite.T(), zlib.Option, suite.builders[ComponentZlib].option())
}

func TestBuilderTestSuite(t *testing.T) {
	suite.Run(t, new(BuilderTestSuite))
}
