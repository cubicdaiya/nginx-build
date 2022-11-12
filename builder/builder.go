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
}

var (
	nginxVersionRe       *regexp.Regexp
	pcreVersionRe        *regexp.Regexp
	zlibVersionRe        *regexp.Regexp
	opensslVersionRe     *regexp.Regexp
	libresslVersionRe    *regexp.Regexp
	opensslQuicVersionRe *regexp.Regexp
	openrestyVersionRe   *regexp.Regexp
	tengineVersionRe     *regexp.Regexp
	nginxQuicVersionRe   *regexp.Regexp
)

func init() {
	nginxVersionRe = regexp.MustCompile(`nginx version: nginx.(\d+\.\d+\.\d+)`)
	pcreVersionRe = regexp.MustCompile(`--with-pcre=.+/pcre-(\d+\.\d+)`)
	zlibVersionRe = regexp.MustCompile(`--with-zlib=.+/zlib-(\d+\.\d+\.\d+)`)
	opensslVersionRe = regexp.MustCompile(`--with-openssl=.+/openssl-(\d+\.\d+\.\d+[a-z]*)`)
	libresslVersionRe = regexp.MustCompile(`--with-openssl=.+/libressl-(\d+\.\d+\.\d+)`)
	opensslQuicVersionRe = regexp.MustCompile(`--with-openssl=.+/openssl-openssl-\d+\.\d+\.\d+-quic1`)
	openrestyVersionRe = regexp.MustCompile(`nginx version: openresty/(\d+\.\d+\.\d+\.\d+)`)
	tengineVersionRe = regexp.MustCompile(`Tengine version: Tengine/(\d+\.\d+\.\d+)`)
	nginxQuicVersionRe = regexp.MustCompile(`nginx version: nginx.(\d+\.\d+\.\d+)`)
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
	case ComponentOpenSSLQuic:
		name = "openssl"
	case ComponentZlib:
		name = "zlib"
	case ComponentOpenResty:
		name = openresty.Name(builder.Version)
	case ComponentTengine:
		name = "tengine"
	case ComponentNginxQuic:
		name = "nginx"
	default:
		panic("invalid component")
	}
	return name
}

func (builder *Builder) option() string {
	name := builder.name()

	// libressl and openssl-quic does not match option name
	if name == "libressl" || name == "openssl-quic" {
		name = "openssl"
	}

	// pcre2 does not match option name
	if name == "pcre2" {
		name = "pcre"
	}

	return fmt.Sprintf("--with-%s", name)
}

func (builder *Builder) DownloadURL() string {
	if builder.Component == ComponentNginxQuic {
		return fmt.Sprintf("%s/%s", builder.DownloadURLPrefix, "tip.tar.gz")
	}
	if builder.Component == ComponentOpenSSLQuic {
		return fmt.Sprintf("%s/openssl-%s+quic1.tar.gz", builder.DownloadURLPrefix, builder.Version)
	}
	return fmt.Sprintf("%s/%s", builder.DownloadURLPrefix, builder.ArchivePath())
}

func (builder *Builder) SourcePath() string {
	if builder.Component == ComponentNginxQuic {
		return "nginx-quic-6cf8ed15fd00"
	}
	if builder.Component == ComponentOpenSSLQuic {
		return fmt.Sprintf("openssl-openssl-%s+quic1", builder.Version)
	}
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
	case "openssl-quic":
		versionRe = opensslQuicVersionRe
	case "tengine":
		versionRe = tengineVersionRe
	case "nginx-quic":
		versionRe = nginxQuicVersionRe
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
		builder.DownloadURLPrefix = OpenSSLDownloadURLPrefix
	case ComponentLibreSSL:
		builder.DownloadURLPrefix = LibreSSLDownloadURLPrefix
	case ComponentOpenSSLQuic:
		builder.DownloadURLPrefix = OpenSSLQuicDownloadURLPrefix
	case ComponentZlib:
		builder.DownloadURLPrefix = ZlibDownloadURLPrefix
	case ComponentOpenResty:
		builder.DownloadURLPrefix = OpenRestyDownloadURLPrefix
	case ComponentTengine:
		builder.DownloadURLPrefix = TengineDownloadURLPrefix
	case ComponentNginxQuic:
		builder.DownloadURLPrefix = NginxQuicDownloadURLPrefix
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
