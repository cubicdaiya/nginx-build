package main

const NGINX_BUILD_VERSION = "0.0.6"

// nginx
const (
	NGINX_VERSION             = "1.7.8"
	NGINX_DOWNLOAD_URL_PREFIX = "http://nginx.org/download"
)

// pcre
const (
	PCRE_VERSION             = "8.36"
	PCRE_DOWNLOAD_URL_PREFIX = "http://ftp.csx.cam.ac.uk/pub/software/programming/pcre"
)

// openssl
const (
	OPENSSL_VERSION             = "1.0.1j"
	OPENSSL_DOWNLOAD_URL_PREFIX = "http://www.openssl.org/source"
)

// zlib
const (
	ZLIB_VERSION             = "1.2.8"
	ZLIB_DOWNLOAD_URL_PREFIX = "http://zlib.net"
)

// component enumerations
const (
	COMPONENT_NGINX = iota
	COMPONENT_PCRE
	COMPONENT_OPENSSL
	COMPONENT_ZLIB
	COMPONENT_MAX
)
