package builder

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/cubicdaiya/nginx-build/command"
	"github.com/cubicdaiya/nginx-build/openresty"
)

const DefaultDownloadTimeout = time.Duration(900) * time.Second

type Builder struct {
	Version           string
	DownloadURLPrefix string
	Component         int
	DownloadTimeout   time.Duration
	// for dependencies such as pcre and zlib and openssl
	Static bool
}

var (
	nginxVersionRe     *regexp.Regexp
	pcreVersionRe      *regexp.Regexp
	zlibVersionRe      *regexp.Regexp
	opensslVersionRe   *regexp.Regexp
	libresslVersionRe  *regexp.Regexp
	openrestyVersionRe *regexp.Regexp
	tengineVersionRe   *regexp.Regexp
)

func init() {
	nginxVersionRe = regexp.MustCompile(`nginx version: nginx.(\d+\.\d+\.\d+)`)
	pcreVersionRe = regexp.MustCompile(`--with-pcre=.+/pcre-(\d+\.\d+)`)
	zlibVersionRe = regexp.MustCompile(`--with-zlib=.+/zlib-(\d+\.\d+\.\d+)`)
	opensslVersionRe = regexp.MustCompile(`--with-openssl=.+/openssl-(\d+\.\d+\.\d+[a-z]*)`)
	libresslVersionRe = regexp.MustCompile(`--with-openssl=.+/libressl-(\d+\.\d+\.\d+)`)
	openrestyVersionRe = regexp.MustCompile(`nginx version: openresty/(\d+\.\d+\.\d+\.\d+)`)
	tengineVersionRe = regexp.MustCompile(`Tengine version: Tengine/(\d+\.\d+\.\d+)`)
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
	name := builder.name()

	// libressl does not match option name
	if name == "libressl" {
		name = "openssl"
	}

	// pcre2 does not match option name
	if name == "pcre2" {
		name = "pcre"
	}

	return fmt.Sprintf("--with-%s", name)
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
	case "tengine":
		versionRe = tengineVersionRe
	}

	m := versionRe.FindSubmatch(result)
	if len(m) < 2 {
		return "", nil
	}
	return string(m[1]), nil
}

func MakeBuilder(component int, version string, timeout time.Duration) Builder {
	var builder Builder
	builder.Component = component
	builder.Version = version
	builder.DownloadTimeout = timeout
	switch component {
	case ComponentNginx:
		builder.DownloadURLPrefix = NginxDownloadURLPrefix
	case ComponentPcre:
		builder.DownloadURLPrefix = fmt.Sprintf("%s/pcre2-%s", PcreDownloadURLPrefix, version)
	case ComponentOpenSSL:
		builder.DownloadURLPrefix = OpenSSLDownloadURLPrefix
	case ComponentLibreSSL:
		builder.DownloadURLPrefix = LibreSSLDownloadURLPrefix
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

func MakeLibraryBuilder(component int, version string, timeout time.Duration, static bool) Builder {
	builder := MakeBuilder(component, version, timeout)
	builder.Static = static
	return builder
}
