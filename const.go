package main

const NGINX_BUILD_VERSION = "0.2.0"

// nginx
const (
	NGINX_VERSION             = "1.7.9"
	NGINX_DOWNLOAD_URL_PREFIX = "http://nginx.org/download"
)

// pcre
const (
	PCRE_VERSION             = "8.36"
	PCRE_DOWNLOAD_URL_PREFIX = "http://ftp.csx.cam.ac.uk/pub/software/programming/pcre"
)

// openssl
const (
	OPENSSL_VERSION             = "1.0.1k"
	OPENSSL_DOWNLOAD_URL_PREFIX = "http://www.openssl.org/source"
)

// zlib
const (
	ZLIB_VERSION             = "1.2.8"
	ZLIB_DOWNLOAD_URL_PREFIX = "http://zlib.net"
)

// openResty
const (
	OPENRESTY_VERSION             = "1.7.7.1"
	OPENRESTY_DOWNLOAD_URL_PREFIX = "http://openresty.org/download"
)

// component enumerations
const (
	COMPONENT_NGINX = iota
	COMPONENT_OPENRESTY
	COMPONENT_PCRE
	COMPONENT_OPENSSL
	COMPONENT_ZLIB
	COMPONENT_MAX
)
