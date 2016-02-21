package main

type ConfigureOptions struct {
	Values map[string]ConfigureOptionValue
	Bools  map[string]ConfigureOptionBool
}

type ConfigureOptionValue struct {
	Name  string
	Desc  string
	Value *string
}

type ConfigureOptionBool struct {
	Name    string
	Desc    string
	Enabled *bool
}

func makeArgsBool() map[string]ConfigureOptionBool {
	argsBool := make(map[string]ConfigureOptionBool)

	argsBool["with-select_module"] = ConfigureOptionBool{Name: "--with-select_module", Desc: "enable select module"}
	argsBool["without-select_module"] = ConfigureOptionBool{Name: "--without-select_module", Desc: "disable select module"}
	argsBool["with-poll_module"] = ConfigureOptionBool{Name: "--with-poll_module", Desc: "enable poll module"}
	argsBool["without-poll_module"] = ConfigureOptionBool{Name: "--without-poll_module", Desc: "disable poll module"}
	argsBool["with-threads"] = ConfigureOptionBool{Name: "--with-threads", Desc: "enable thread pool support"}
	argsBool["with-file-aio"] = ConfigureOptionBool{Name: "--with-file-aio", Desc: "enable file AIO support"}
	argsBool["with-ipv6"] = ConfigureOptionBool{Name: "--with-ipv6", Desc: "enable IPv6 support"}
	argsBool["with-http_ssl_module"] = ConfigureOptionBool{Name: "--with-http_ssl_module", Desc: "enable ngx_http_ssl_module"}
	argsBool["with-http_spdy_module"] = ConfigureOptionBool{Name: "--with-http_spdy_module", Desc: "enable ngx_http_spdy_module"}
	argsBool["with-http_v2_module"] = ConfigureOptionBool{Name: "--with-http_v2_module", Desc: "enable ngx_http_v2_module"}
	argsBool["with-http_realip_module"] = ConfigureOptionBool{Name: "--with-http_realip_module", Desc: "enable ngx_http_realip_module"}
	argsBool["with-http_addition_module"] = ConfigureOptionBool{Name: "--with-http_addition_module", Desc: "enable ngx_http_addition_module"}
	argsBool["with-http_xslt_module"] = ConfigureOptionBool{Name: "--with-http_xslt_module", Desc: "enable ngx_http_xslt_module"}
	argsBool["with-http_image_filter_module"] = ConfigureOptionBool{Name: "--with-http_image_filter_module", Desc: "enable ngx_http_image_filter_module"}
	argsBool["with-http_geoip_module"] = ConfigureOptionBool{Name: "--with-http_geoip_module", Desc: "enable ngx_http_geoip_module"}
	argsBool["with-http_sub_module"] = ConfigureOptionBool{Name: "--with-http_sub_module", Desc: "enable ngx_http_sub_module"}
	argsBool["with-http_dav_module"] = ConfigureOptionBool{Name: "--with-http_dav_module", Desc: "enable ngx_http_dav_module"}
	argsBool["with-http_flv_module"] = ConfigureOptionBool{Name: "--with-http_flv_module", Desc: "enable ngx_http_flv_module"}
	argsBool["with-http_mp4_module"] = ConfigureOptionBool{Name: "--with-http_mp4_module", Desc: "enable ngx_http_mp4_module"}
	argsBool["with-http_gunzip_module"] = ConfigureOptionBool{Name: "--with-http_gunzip_module", Desc: "enable ngx_http_gunzip_module"}
	argsBool["with-http_gzip_static_module"] = ConfigureOptionBool{Name: "--with-http_gzip_static_module", Desc: "enable ngx_http_gzip_static_module"}
	argsBool["with-http_auth_request_module"] = ConfigureOptionBool{Name: "--with-http_auth_request_module", Desc: "enable ngx_http_auth_request_module"}
	argsBool["with-http_random_index_module"] = ConfigureOptionBool{Name: "--with-http_random_index_module", Desc: "enable ngx_http_random_index_module"}
	argsBool["with-http_secure_link_module"] = ConfigureOptionBool{Name: "--with-http_secure_link_module", Desc: "enable ngx_http_secure_link_module"}
	argsBool["with-http_degradation_module"] = ConfigureOptionBool{Name: "--with-http_degradation_module", Desc: "enable ngx_http_degradation_module"}
	argsBool["with-http_stub_status_module"] = ConfigureOptionBool{Name: "--with-http_stub_status_module", Desc: "enable ngx_http_stub_status_module"}
	argsBool["without-http_charset_module"] = ConfigureOptionBool{Name: "--with-http_charset_module", Desc: "disable ngx_http_charset_module"}
	argsBool["without-http_gzip_module"] = ConfigureOptionBool{Name: "--with-http_gzip_module", Desc: "disable ngx_http_gzip_module"}
	argsBool["without-http_ssi_module"] = ConfigureOptionBool{Name: "--with-http_ssi_module", Desc: "disable ngx_http_ssi_module"}
	argsBool["without-http_userid_module"] = ConfigureOptionBool{Name: "--with-http_userid_module", Desc: "disable ngx_http_userid_module"}
	argsBool["without-http_access_module"] = ConfigureOptionBool{Name: "--with-http_access_module", Desc: "disable ngx_http_access_module"}
	argsBool["without-http_auth_basic_module"] = ConfigureOptionBool{Name: "--with-http_auth_basic_module", Desc: "disable ngx_http_auth_basic_module"}
	argsBool["without-http_autoindex_module"] = ConfigureOptionBool{Name: "--with-http_autoindex_module", Desc: "disable ngx_http_autoindex_module"}
	argsBool["without-http_geo_module"] = ConfigureOptionBool{Name: "--with-http_geo_module", Desc: "disable ngx_http_geo_module"}
	argsBool["without-http_map_module"] = ConfigureOptionBool{Name: "--with-http_map_module", Desc: "disable ngx_http_map_module"}
	argsBool["without-http_split_clients_module"] = ConfigureOptionBool{Name: "--with-http_split_clients_module", Desc: "disable ngx_http_split_clients_module"}
	argsBool["without-http_referer_module"] = ConfigureOptionBool{Name: "--with-http_referer_module", Desc: "disable ngx_http_referer_module"}
	argsBool["without-http_rewrite_module"] = ConfigureOptionBool{Name: "--with-http_rewrite_module", Desc: "disable ngx_http_rewrite_module"}
	argsBool["without-http_proxy_module"] = ConfigureOptionBool{Name: "--with-http_proxy_module", Desc: "disable ngx_http_proxy_module"}
	argsBool["without-http_fastcgi_module"] = ConfigureOptionBool{Name: "--with-http_fastcgi_module", Desc: "disable ngx_http_fastcgi_module"}
	argsBool["without-http_uwsgi_module"] = ConfigureOptionBool{Name: "--with-http_uwsgi_module", Desc: "disable ngx_http_uwsgi_module"}
	argsBool["without-http_scgi_module"] = ConfigureOptionBool{Name: "--with-http_scgi_module", Desc: "disable ngx_http_scgi_module"}
	argsBool["without-http_memcached_module"] = ConfigureOptionBool{Name: "--with-http_memcached_module", Desc: "disable ngx_http_memcached_module"}
	argsBool["without-http_limit_conn_module"] = ConfigureOptionBool{Name: "--with-http_limit_conn_module", Desc: "disable ngx_http_limit_conn_module"}
	argsBool["without-http_limit_req_module"] = ConfigureOptionBool{Name: "--with-http_limit_req_module", Desc: "disable ngx_http_limit_req_module"}
	argsBool["without-http_empty_gif_module"] = ConfigureOptionBool{Name: "--with-http_empty_gif_module", Desc: "disable ngx_http_empty_gif_module"}
	argsBool["without-http_browser_module"] = ConfigureOptionBool{Name: "--with-http_browser_module", Desc: "disable ngx_http_browser_module"}
	argsBool["without-http_upstream_hash_module"] = ConfigureOptionBool{Name: "--with-http_upstream_hash_module", Desc: "disable ngx_http_upstream_hash_module"}
	argsBool["without-http_upstream_ip_hash_module"] = ConfigureOptionBool{Name: "--with-http_upstream_ip_hash_module", Desc: "disable ngx_http_upstream_ip_hash_module"}
	argsBool["without-http_upstream_least_conn_module"] = ConfigureOptionBool{Name: "--with-http_upstream_least_conn_module", Desc: "disable ngx_http_upstream_least_conn_module"}
	argsBool["without-http_upstream_keepalive_module"] = ConfigureOptionBool{Name: "--with-http_upstream_keepalive_module", Desc: "disable ngx_http_upstream_keepalive_module"}
	argsBool["without-http_upstream_zone_module"] = ConfigureOptionBool{Name: "--with-http_upstream_zone_module", Desc: "disable ngx_http_upstream_zone_module"}
	argsBool["with-http_perl_module"] = ConfigureOptionBool{Name: "--with-http_perl_module", Desc: "enable ngx_http_perl_module"}
	argsBool["without-http"] = ConfigureOptionBool{Name: "--without-http", Desc: "disable HTTP server"}
	argsBool["without-http-cache"] = ConfigureOptionBool{Name: "--without-http-cache", Desc: "disable HTTP cache"}
	argsBool["with-mail"] = ConfigureOptionBool{Name: "--with-mail", Desc: "enable POP3/IMAP4/SMTP proxy module"}
	argsBool["with-mail_ssl_module"] = ConfigureOptionBool{Name: "--with-mail_ssl_module", Desc: "enable ngx_mail_ssl_module"}
	argsBool["without-mail_pop3_module"] = ConfigureOptionBool{Name: "--without-mail_pop3_module", Desc: "disable ngx_mail_pop3_module"}
	argsBool["without-mail_imap_module"] = ConfigureOptionBool{Name: "--without-mail_imap_module", Desc: "disable ngx_mail_imap_module"}
	argsBool["without-mail_smtp_module"] = ConfigureOptionBool{Name: "--without-mail_smtp_module", Desc: "disable ngx_mail_smtp_module"}
	argsBool["with-stream"] = ConfigureOptionBool{Name: "--with-stream", Desc: "enable TCP proxy module"}
	argsBool["with-stream_ssl_module"] = ConfigureOptionBool{Name: "--with-stream_ssl_module", Desc: "enable ngx_stream_ssl_module"}
	argsBool["without-stream_access_module"] = ConfigureOptionBool{Name: "--without-stream_access_module", Desc: "disable ngx_stream_access_module"}
	argsBool["without-stream_upstream_hash_module"] = ConfigureOptionBool{Name: "--without-stream_upstream_hash_module", Desc: "disable ngx_stream_upstream_hash_module"}
	argsBool["without-stream_upstream_least_conn_module"] = ConfigureOptionBool{Name: "--without-stream_upstream_least_conn_module", Desc: "disable ngx_stream_upstream_least_conn_module"}
	argsBool["without-stream_upstream_zone_module"] = ConfigureOptionBool{Name: "--without-stream_upstream_zone_module", Desc: "disable ngx_stream_upstream_zone_module"}
	argsBool["with-google_perftools_module"] = ConfigureOptionBool{Name: "--with-google_perftools_module", Desc: "enable ngx_google_perftools_module"}
	argsBool["with-cpp_test_module"] = ConfigureOptionBool{Name: "--with-cpp_test_module", Desc: "enable ngx_cpp_test_module"}
	argsBool["without-pcre"] = ConfigureOptionBool{Name: "--without-pcre", Desc: "disable PCRE library usage"}
	argsBool["with-pcre-jit"] = ConfigureOptionBool{Name: "--with-pcre-jit", Desc: "build PCRE with JIT compilation support"}
	argsBool["with-md5-asm"] = ConfigureOptionBool{Name: "--with-md5-asm", Desc: "use md5 assembler sources"}
	argsBool["with-sha1-asm"] = ConfigureOptionBool{Name: "--with-sha1-asm", Desc: "use sha1 assembler sources"}
	argsBool["with-debug"] = ConfigureOptionBool{Name: "--with-debug", Desc: "enable debug logging"}

	//
	// TODO: The options below are not suitable for flag package.
	//
	// argsBool["with-http_xslt_module=dynamic"] = ConfigureOptionBool{Name: "--with-http_xslt_module=dynamic", Desc: "enable dynamic ngx_http_xslt_module"}
	// argsBool["with-http_image_filter_module=dynamic"] = ConfigureOptionBool{Name: "--with-http_image_filter_module=dynamic", Desc: "enable dynamic ngx_http_image_filter_module"}
	// argsBool["with-http_geoip_module=dynamic"] = ConfigureOptionBool{Name: "--with-http_geoip_module=dynamic", Desc: "enable dynamic ngx_http_geoip_module"}
	// argsBool["with-mail=dynamic"] = ConfigureOptionBool{Name: "--with-mail=dynamic", Desc: "enable dynamic POP3/IMAP4/SMTP proxy module"}
	// argsBool["with-stream=dynamic"] = ConfigureOptionBool{Name: "--with-stream=dynamic", Desc: "enable dynamic TCP proxy module"}

	return argsBool
}

func makeArgsString() map[string]ConfigureOptionValue {
	argsString := make(map[string]ConfigureOptionValue)

	argsString["prefix"] = ConfigureOptionValue{Name: "--prefix", Desc: "set installation prefix"}
	argsString["sbin-path"] = ConfigureOptionValue{Name: "--sbin-path", Desc: "set nginx binary pathname"}
	argsString["modules-path"] = ConfigureOptionValue{Name: "--modules-path", Desc: "set modules path"}
	argsString["conf-path"] = ConfigureOptionValue{Name: "--conf-path", Desc: "set ginx.conf pathname"}
	argsString["error-log-path"] = ConfigureOptionValue{Name: "--error-log-path", Desc: "set error log pathname"}
	argsString["pid-path"] = ConfigureOptionValue{Name: "--pid-path", Desc: "set nginx.pid pathname"}
	argsString["lock-path"] = ConfigureOptionValue{Name: "--lock-path", Desc: "set nginx.lock pathname"}
	argsString["user"] = ConfigureOptionValue{Name: "--user", Desc: "set non-privileged user for worker processes"}
	argsString["group"] = ConfigureOptionValue{Name: "--group", Desc: "set non-privileged group for worker processes"}
	argsString["build"] = ConfigureOptionValue{Name: "--build", Desc: "set build name"}
	argsString["builddir"] = ConfigureOptionValue{Name: "--builddir", Desc: "set build directory"}
	argsString["with-perl_modules_path="] = ConfigureOptionValue{Name: "--with-perl_modules_path", Desc: "set Perl modules path"}
	argsString["with-perl"] = ConfigureOptionValue{Name: "--with-perl", Desc: "set perl binary pathname"}
	argsString["http-log-path"] = ConfigureOptionValue{Name: "--http-log-path", Desc: "set http access log pathname"}
	argsString["http-client-body-temp-path"] = ConfigureOptionValue{Name: "--http-client-body-temp-path", Desc: "set path to store http client request body temporary files"}
	argsString["http-proxy-temp-path"] = ConfigureOptionValue{Name: "--http-proxy-temp-path", Desc: "set path to store http proxy temporary files"}
	argsString["http-fastcgi-temp-path"] = ConfigureOptionValue{Name: "--http-fastcgi-temp-path", Desc: "set path to store http fastcgi temporary files"}
	argsString["http-uwsgi-temp-path"] = ConfigureOptionValue{Name: "--http-uwsgi-temp-path", Desc: "set path to store http uwsgi temporary files"}
	argsString["http-scgi-temp-path"] = ConfigureOptionValue{Name: "--http-scgi-temp-path", Desc: "set path to store http scgi temporary files"}
	argsString["add-module"] = ConfigureOptionValue{Name: "--add-module", Desc: "enable external module"}
	argsString["add-dynamic-module"] = ConfigureOptionValue{Name: "--add-dynamic module", Desc: "enable dynamic external module"}
	argsString["with-cc"] = ConfigureOptionValue{Name: "--with-cc", Desc: "set C compiler pathname"}
	argsString["with-cpp"] = ConfigureOptionValue{Name: "--with-cpp", Desc: "set c preprocessor pathname"}
	argsString["with-cc-opt"] = ConfigureOptionValue{Name: "--with-cc-opt", Desc: "set additional C compiler options"}
	argsString["with-ld-opt"] = ConfigureOptionValue{Name: "--with-ld-opt", Desc: "set additional linker options"}
	argsString["with-cpu-opt"] = ConfigureOptionValue{Name: "--with-cpu-opt", Desc: "build for the specified CPU, valid values: pentium, pentiumpro, pentium3, pentium4, athlon, opteron, sparc32, sparc64, ppc64"}
	argsString["with-pcre"] = ConfigureOptionValue{Name: "--with-pcre", Desc: "set path to PCRE library sources"}
	argsString["with-pcre-opt"] = ConfigureOptionValue{Name: "--with-pcre-opt", Desc: "set additional build options for PCRE"}
	argsString["with-md5"] = ConfigureOptionValue{Name: "--with-md5", Desc: "set path to md5 library sources"}
	argsString["with-md5-opt"] = ConfigureOptionValue{Name: "--with-md5-opt", Desc: "set additional build options for md5"}
	argsString["with-sha1"] = ConfigureOptionValue{Name: "--with-sha1", Desc: "set path to sha1 library sources"}
	argsString["with-sha1-opt"] = ConfigureOptionValue{Name: "--with-sha1-opt", Desc: "set additional build options for sha1"}
	argsString["with-zlib"] = ConfigureOptionValue{Name: "--with-zlib", Desc: "set path to zlib library sources"}
	argsString["with-zlib-opt"] = ConfigureOptionValue{Name: "--with-zlib-opt", Desc: "set additional build options for zlib"}
	argsString["with-zlib-asm"] = ConfigureOptionValue{Name: "--with-zlib-asm", Desc: "use zlib assembler sources optimized for the specified CPU, valid values: pentium, pentiumpro"}
	argsString["with-libatomic"] = ConfigureOptionValue{Name: "--with-libatomic", Desc: "set path to libatomic_ops library sources"}
	argsString["with-openssl"] = ConfigureOptionValue{Name: "--with-openssl", Desc: "set path to OpenSSL library sources"}
	argsString["with-openssl-opt"] = ConfigureOptionValue{Name: "--with-openssl-opt", Desc: "set additional build options for OpenSSL"}

	return argsString
}
