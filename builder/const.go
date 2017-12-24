package builder

// nginx
const (
	NginxVersion           = "1.13.7"
	NginxDownloadURLPrefix = "https://nginx.org/download"
)

// pcre
const (
	PcreVersion           = "8.41"
	PcreDownloadURLPrefix = "https://ftp.pcre.org/pub/pcre"
)

// openssl
const (
	OpenSSLVersion           = "1.0.2n"
	OpenSSLDownloadURLPrefix = "https://www.openssl.org/source"
)

// libressl
const (
	LibreSSLVersion           = "2.6.3"
	LibreSSLDownloadURLPrefix = "https://ftp.openbsd.org/pub/OpenBSD/LibreSSL"
)

// zlib
const (
	ZlibVersion           = "1.2.11"
	ZlibDownloadURLPrefix = "https://zlib.net/fossils"
)

// openResty
const (
	OpenRestyVersion           = "1.13.6.1"
	OpenRestyDownloadURLPrefix = "https://openresty.org/download"
)

// tengine
const (
	TengineVersion           = "2.2.1"
	TengineDownloadURLPrefix = "http://tengine.taobao.org/download"
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
