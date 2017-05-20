package builder

// nginx
const (
	NginxVersion           = "1.13.0"
	NginxDownloadURLPrefix = "https://nginx.org/download"
)

// pcre
const (
	PcreVersion           = "8.40"
	PcreDownloadURLPrefix = "http://ftp.csx.cam.ac.uk/pub/software/programming/pcre"
)

// openssl
const (
	OpenSSLVersion           = "1.0.2k"
	OpenSSLDownloadURLPrefix = "https://www.openssl.org/source"
)

// zlib
const (
	ZlibVersion           = "1.2.11"
	ZlibDownloadURLPrefix = "https://zlib.net/fossils"
)

// openResty
const (
	OpenRestyVersion           = "1.11.2.3"
	OpenRestyDownloadURLPrefix = "https://openresty.org/download"
)

// tengine
const (
	TengineVersion           = "2.2.0"
	TengineDownloadURLPrefix = "http://tengine.taobao.org/download"
)

// component enumerations
const (
	ComponentNginx = iota
	ComponentOpenResty
	ComponentTengine
	ComponentPcre
	ComponentOpenSSL
	ComponentZlib
	ComponentMax
)
