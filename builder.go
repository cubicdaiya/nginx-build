package main

import (
	"fmt"
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
	}
	return name
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
	}
	return builder
}
