package main

import (
	"fmt"
	"strings"
)

type Builder struct {
	Version           string
	DownloadURLPrefix string
	Component         int
}

func (builder *Builder) name() string {
	var name string
	switch builder.Component {
	case COMPONENT_NGINX:
		name = "nginx"
	case COMPONENT_PCRE:
		name = "pcre"
	case COMPONENT_OPENSSL:
		name = "openssl"
	case COMPONENT_ZLIB:
		name = "zlib"
	case COMPONENT_OPENRESTY:
		name = "ngx_openresty"
	default:
		panic("invalid component")
	}
	return name
}

func (builder *Builder) option() string {
	return fmt.Sprintf("--with-%s", builder.name())
}

func (builder *Builder) downloadURL() string {
	return fmt.Sprintf("%s/%s", builder.DownloadURLPrefix, builder.archivePath())
}

func (builder *Builder) sourcePath() string {
	return fmt.Sprintf("%s-%s", builder.name(), builder.Version)
}

func (builder *Builder) archivePath() string {
	return fmt.Sprintf("%s.tar.gz", builder.sourcePath())
}

func (builder *Builder) isIncludeWithOption(nginxConfigure string) bool {
	if strings.Contains(nginxConfigure, builder.option()+"=") {
		return true
	}
	return false
}

func (builder *Builder) warnMsgWithLibrary() string {
	return fmt.Sprintf("[warn]Using '%s' is discouraged. Instead give '-%s' and '-%sversion' to 'nginx-build'",
		builder.option(), builder.name(), builder.name())
}

func makeBuilder(component int, version string) Builder {
	var builder Builder
	builder.Component = component
	builder.Version = version
	switch component {
	case COMPONENT_NGINX:
		builder.DownloadURLPrefix = NGINX_DOWNLOAD_URL_PREFIX
	case COMPONENT_PCRE:
		builder.DownloadURLPrefix = PCRE_DOWNLOAD_URL_PREFIX
	case COMPONENT_OPENSSL:
		builder.DownloadURLPrefix = OPENSSL_DOWNLOAD_URL_PREFIX
	case COMPONENT_ZLIB:
		builder.DownloadURLPrefix = ZLIB_DOWNLOAD_URL_PREFIX
	case COMPONENT_OPENRESTY:
		builder.DownloadURLPrefix = OPENRESTY_DOWNLOAD_URL_PREFIX
	default:
		panic("invalid component")
	}
	return builder
}
