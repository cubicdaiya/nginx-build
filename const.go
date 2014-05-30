package main

const NGINX_BUILD_VERSION = "0.0.1"

const NGINX_DOWNLOAD_URL_PREFIX = "http://nginx.org/download"
const NGINX_VERSION = "1.7.0"

const PCRE_DOWNLOAD_URL_PREFIX = "http://ftp.csx.cam.ac.uk/pub/software/programming/pcre"
const PCRE_VERSION = "8.35"

const OPENSSL_DOWNLOAD_URL_PREFIX = "http://www.openssl.org/source"
const OPENSSL_VERSION = "1.0.1g"

const ZLIB_DOWNLOAD_URL_PREFIX = "http://zlib.net"
const ZLIB_VERSION = "1.2.8"

const (
	COMPONENT_NGINX = iota
	COMPONENT_PCRE
	COMPONENT_OPENSSL
	COMPONENT_ZLIB
	COMPONENT_MAX
)
