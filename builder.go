package main

type Builder struct {
	Version string
	DownLoadPrefix string
	Component int
	WorkDir string
}

func MakeBuilder(component int, version string) Builder {
	var builder Builder
	builder.Component = component
	builder.Version = version
	switch (component) {
		case COMPONENT_NGINX:
		builder.DownLoadPrefix = NGINX_DOWNLOAD_URL_PREFIX
		case COMPONENT_PCRE:
		builder.DownLoadPrefix = PCRE_DOWNLOAD_URL_PREFIX
		case COMPONENT_OPENSSL:
		builder.DownLoadPrefix = OPENSSL_DOWNLOAD_URL_PREFIX
	}
	return builder
}
