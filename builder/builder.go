package builder

import (
	"fmt"
	"strings"

	"github.com/cubicdaiya/nginx-build/openresty"
)

type Builder struct {
	Version           string
	DownloadURLPrefix string
	Component         int
}

func (builder *Builder) name() string {
	var name string
	switch builder.Component {
	case ComponentNginx:
		name = "nginx"
	case ComponentPcre:
		name = "pcre"
	case ComponentOpenSSL:
		name = "openssl"
	case ComponentZlib:
		name = "zlib"
	case ComponentOpenResty:
		name = openresty.Name(builder.Version)
	case ComponentTengine:
		name = "tengine"
	default:
		panic("invalid component")
	}
	return name
}

func (builder *Builder) option() string {
	return fmt.Sprintf("--with-%s", builder.name())
}

func (builder *Builder) DownloadURL() string {
	return fmt.Sprintf("%s/%s", builder.DownloadURLPrefix, builder.ArchivePath())
}

func (builder *Builder) SourcePath() string {
	return fmt.Sprintf("%s-%s", builder.name(), builder.Version)
}

func (builder *Builder) ArchivePath() string {
	return fmt.Sprintf("%s.tar.gz", builder.SourcePath())
}

func (builder *Builder) LogPath() string {
	return fmt.Sprintf("%s-%s.log", builder.name(), builder.Version)
}

func (builder *Builder) IsIncludeWithOption(nginxConfigure string) bool {
	if strings.Contains(nginxConfigure, builder.option()+"=") {
		return true
	}
	return false
}

func (builder *Builder) WarnMsgWithLibrary() string {
	return fmt.Sprintf("[warn]Using '%s' is discouraged. Instead give '-%s' and '-%sversion' to 'nginx-build'",
		builder.option(), builder.name(), builder.name())
}

func MakeBuilder(component int, version string) Builder {
	var builder Builder
	builder.Component = component
	builder.Version = version
	switch component {
	case ComponentNginx:
		builder.DownloadURLPrefix = NginxDownloadURLPrefix
	case ComponentPcre:
		builder.DownloadURLPrefix = PcreDownloadURLPrefix
	case ComponentOpenSSL:
		builder.DownloadURLPrefix = OpenSSLDownloadURLPrefix
	case ComponentZlib:
		builder.DownloadURLPrefix = ZlibDownloadURLPrefix
	case ComponentOpenResty:
		builder.DownloadURLPrefix = OpenRestyDownloadURLPrefix
	case ComponentTengine:
		builder.DownloadURLPrefix = TengineDownloadURLPrefix
	default:
		panic("invalid component")
	}
	return builder
}
