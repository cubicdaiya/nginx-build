package builder

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/cubicdaiya/nginx-build/command"
	"github.com/cubicdaiya/nginx-build/openresty"
)

type Builder struct {
	Version           string
	DownloadURLPrefix string
	Component         int
	// for dependencies such as pcre and zlib and openssl
	Static bool
	// for custom SSL library
	CustomURL  string
	CustomName string
	CustomTag  string
}

var (
	nginxVersionRe     *regexp.Regexp
	pcreVersionRe      *regexp.Regexp
	zlibVersionRe      *regexp.Regexp
	opensslVersionRe   *regexp.Regexp
	libresslVersionRe  *regexp.Regexp
	openrestyVersionRe *regexp.Regexp
	freenginxVersionRe *regexp.Regexp
)

func init() {
	nginxVersionRe = regexp.MustCompile(`nginx version: nginx.(\d+\.\d+\.\d+)`)
	pcreVersionRe = regexp.MustCompile(`--with-pcre=.+/pcre-(\d+\.\d+)`)
	zlibVersionRe = regexp.MustCompile(`--with-zlib=.+/zlib-(\d+\.\d+\.\d+)`)
	opensslVersionRe = regexp.MustCompile(`--with-openssl=.+/openssl-(\d+\.\d+\.\d+[a-z]*)`)
	libresslVersionRe = regexp.MustCompile(`--with-openssl=.+/libressl-(\d+\.\d+\.\d+)`)
	openrestyVersionRe = regexp.MustCompile(`nginx version: openresty/(\d+\.\d+\.\d+\.\d+)`)
	freenginxVersionRe = regexp.MustCompile(`freenginx version: freenginx/(\d+\.\d+\.\d+)`)
}

func (builder *Builder) name() string {
	var name string
	switch builder.Component {
	case ComponentNginx:
		name = "nginx"
	case ComponentPcre:
		name = "pcre2"
	case ComponentOpenSSL:
		name = "openssl"
	case ComponentLibreSSL:
		name = "libressl"
	case ComponentCustomSSL:
		if builder.CustomName != "" {
			name = builder.CustomName
		} else {
			name = "customssl"
		}
	case ComponentZlib:
		name = "zlib"
	case ComponentOpenResty:
		name = openresty.Name(builder.Version)
	case ComponentFreenginx:
		name = "freenginx"
	default:
		panic("invalid component")
	}
	return name
}

func (builder *Builder) option() string {
	name := builder.name()

	// libressl does not match option name
	if name == "libressl" {
		name = "openssl"
	}

	// custom ssl defaults to openssl option
	if builder.Component == ComponentCustomSSL {
		name = "openssl"
	}

	// pcre2 does not match option name
	if name == "pcre2" {
		name = "pcre"
	}

	return fmt.Sprintf("--with-%s", name)
}

func (builder *Builder) DownloadURL() string {
	switch builder.Component {
	case ComponentNginx:
		return fmt.Sprintf("%s/nginx-%s.tar.gz", NginxDownloadURLPrefix, builder.Version)
	case ComponentPcre:
		return fmt.Sprintf("%s/pcre2-%s/pcre2-%s.tar.gz", PcreDownloadURLPrefix, builder.Version, builder.Version)
	case ComponentOpenSSL:
		return fmt.Sprintf("%s/openssl-%s/openssl-%s.tar.gz", OpenSSLDownloadURLPrefix, builder.Version, builder.Version)
	case ComponentLibreSSL:
		return fmt.Sprintf("%s/libressl-%s.tar.gz", LibreSSLDownloadURLPrefix, builder.Version)
	case ComponentCustomSSL:
		return builder.CustomURL
	case ComponentZlib:
		return fmt.Sprintf("%s/zlib-%s.tar.gz", ZlibDownloadURLPrefix, builder.Version)
	case ComponentOpenResty:
		return fmt.Sprintf("%s/openresty-%s.tar.gz", OpenRestyDownloadURLPrefix, builder.Version)
	case ComponentFreenginx:
		return fmt.Sprintf("%s/freenginx-%s.tar.gz", FreenginxDownloadURLPrefix, builder.Version)
	default:
		panic("invalid component")
	}
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
	return strings.Contains(nginxConfigure, builder.option()+"=")
}

func (builder *Builder) WarnMsgWithLibrary() string {
	return fmt.Sprintf("[warn]Using '%s' is discouraged. Instead give '-%s' and '-%sversion' to 'nginx-build'",
		builder.option(), builder.name(), builder.name())
}

func (builder *Builder) InstalledVersion() (string, error) {
	nginxBinPath := "/usr/local/sbin/nginx"
	if os.Getenv("NGINX_BIN") != "" {
		nginxBinPath = os.Getenv("NGINX_BIN")
	}
	args := []string{nginxBinPath, "-V"}
	cmd, err := command.Make(args)
	if err != nil {
		return "", err
	}

	result, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	openRestyName := openresty.Name(builder.Version)
	var versionRe *regexp.Regexp

	switch builder.name() {
	case "nginx":
		versionRe = nginxVersionRe
	case openRestyName:
		versionRe = openrestyVersionRe
	case "zlib":
		versionRe = zlibVersionRe
	case "pcre2":
		versionRe = pcreVersionRe
	case "openssl":
		versionRe = opensslVersionRe
	case "libressl":
		versionRe = libresslVersionRe
	case "freenginx":
		versionRe = freenginxVersionRe
	}

	m := versionRe.FindSubmatch(result)
	if len(m) < 2 {
		return "", nil
	}
	return string(m[1]), nil
}

func MakeBuilder(component int, version string) Builder {
	var builder Builder
	builder.Component = component
	builder.Version = version
	switch component {
	case ComponentNginx:
		builder.DownloadURLPrefix = NginxDownloadURLPrefix
	case ComponentPcre:
		builder.DownloadURLPrefix = fmt.Sprintf("%s/pcre2-%s", PcreDownloadURLPrefix, version)
	case ComponentOpenSSL:
		builder.DownloadURLPrefix = fmt.Sprintf("%s/openssl-%s", OpenSSLDownloadURLPrefix, version)
	case ComponentLibreSSL:
		builder.DownloadURLPrefix = LibreSSLDownloadURLPrefix
	case ComponentZlib:
		builder.DownloadURLPrefix = ZlibDownloadURLPrefix
	case ComponentOpenResty:
		builder.DownloadURLPrefix = OpenRestyDownloadURLPrefix
	case ComponentFreenginx:
		builder.DownloadURLPrefix = FreenginxDownloadURLPrefix
	default:
		panic("invalid component")
	}
	return builder
}

func MakeLibraryBuilder(component int, version string, static bool) Builder {
	builder := MakeBuilder(component, version)
	builder.Static = static
	return builder
}
