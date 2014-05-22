package main

import (
	"fmt"
)

func (builder *Builder) SourcePath() string {
	var result string
	switch (builder.Component) {
	case COMPONENT_NGINX:
		result = fmt.Sprintf("%s-%s", "nginx", builder.Version)
	case COMPONENT_PCRE:
		result = fmt.Sprintf("%s-%s", "pcre", builder.Version)
	case COMPONENT_OPENSSL:
		result = fmt.Sprintf("%s-%s", "openssl", builder.Version)
	}
	return result
}

func (builder *Builder) ArchivePath() string {
	return fmt.Sprintf("%s.tar.gz", builder.SourcePath())
}
