package builder

// nginx
const (
	NginxVersion           = "1.26.1"
	NginxDownloadURLPrefix = "https://nginx.org/download"
)

// pcre
const (
	PcreVersion           = "10.44"
	PcreDownloadURLPrefix = "https://github.com/PCRE2Project/pcre2/releases/download"
)

// openssl
const (
	OpenSSLVersion           = "3.3.1"
	OpenSSLDownloadURLPrefix = "https://www.openssl.org/source"
)

// libressl
const (
	LibreSSLVersion           = "3.9.2"
	LibreSSLDownloadURLPrefix = "https://ftp.openbsd.org/pub/OpenBSD/LibreSSL"
)

// zlib
const (
	ZlibVersion           = "1.3.1"
	ZlibDownloadURLPrefix = "https://zlib.net"
)

// openResty
const (
	OpenRestyVersion           = "1.25.3.1"
	OpenRestyDownloadURLPrefix = "https://openresty.org/download"
)

// freenginx
const (
	FreenginxVersion           = "1.26.0"
	FreenginxDownloadURLPrefix = "https://freenginx.org/download"
)

// component enumerations
const (
	ComponentNginx = iota
	ComponentOpenResty
	ComponentFreenginx
	ComponentPcre
	ComponentOpenSSL
	ComponentLibreSSL
	ComponentZlib
	ComponentMax
)
