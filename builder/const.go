package builder

// nginx
const (
	NginxVersion           = "1.19.2"
	NginxDownloadURLPrefix = "https://nginx.org/download"
)

// pcre
const (
	PcreVersion           = "8.44"
	PcreDownloadURLPrefix = "https://ftp.pcre.org/pub/pcre"
)

// openssl
const (
	OpenSSLVersion           = "1.1.1g"
	OpenSSLDownloadURLPrefix = "https://www.openssl.org/source"
)

// libressl
const (
	LibreSSLVersion           = "3.0.2"
	LibreSSLDownloadURLPrefix = "https://ftp.openbsd.org/pub/OpenBSD/LibreSSL"
)

// zlib
const (
	ZlibVersion           = "1.2.11"
	ZlibDownloadURLPrefix = "https://zlib.net/fossils"
)

// openResty
const (
	OpenRestyVersion           = "1.17.8.2"
	OpenRestyDownloadURLPrefix = "https://openresty.org/download"
)

// tengine
const (
	TengineVersion           = "2.3.2"
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
