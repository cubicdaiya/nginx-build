package builder

// nginx
const (
	NginxVersion           = "1.21.5"
	NginxDownloadURLPrefix = "https://nginx.org/download"
)

// pcre
const (
	PcreVersion           = "10.39"
	PcreDownloadURLPrefix = "https://github.com/PhilipHazel/pcre2/releases/download"
)

// openssl
const (
	OpenSSLVersion           = "1.1.1m"
	OpenSSLDownloadURLPrefix = "https://www.openssl.org/source"
)

// libressl
const (
	LibreSSLVersion           = "3.4.2"
	LibreSSLDownloadURLPrefix = "https://ftp.openbsd.org/pub/OpenBSD/LibreSSL"
)

// zlib
const (
	ZlibVersion           = "1.2.11"
	ZlibDownloadURLPrefix = "https://zlib.net/fossils"
)

// openResty
const (
	OpenRestyVersion           = "1.19.9.1"
	OpenRestyDownloadURLPrefix = "https://openresty.org/download"
)

// tengine
const (
	TengineVersion           = "2.3.3"
	TengineDownloadURLPrefix = "https://tengine.taobao.org/download"
)

// component enumerations
const (
	ComponentNginx = iota
	ComponentOpenResty
	ComponentTengine
	ComponentPcre
	ComponentOpenSSL
	ComponentLibreSSL
	ComponentZlib
	ComponentMax
)
