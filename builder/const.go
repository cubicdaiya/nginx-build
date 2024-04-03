package builder

// nginx
const (
	NginxVersion           = "1.24.0"
	NginxDownloadURLPrefix = "https://nginx.org/download"
)

// pcre
const (
	PcreVersion           = "10.43"
	PcreDownloadURLPrefix = "https://github.com/PCRE2Project/pcre2/releases/download"
)

// openssl
const (
	OpenSSLVersion           = "3.2.1"
	OpenSSLDownloadURLPrefix = "https://www.openssl.org/source"
)

// libressl
const (
	LibreSSLVersion           = "3.9.1"
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
	FreenginxVersion           = "1.25.5"
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
