package main

import (
	"fmt"
	"log"

	"github.com/cubicdaiya/nginx-build/builder"
)

func versionsGenNginx() []string {
	return []string{
		fmt.Sprintf("nginx-%s", builder.NginxVersion),
	}
}

func versionsGenOpenResty() []string {
	return []string{
		fmt.Sprintf("openresty-%s", builder.OpenRestyVersion),
	}
}

func versionsGenFreenginx() []string {
	return []string{
		fmt.Sprintf("freenginx-%s", builder.FreenginxVersion),
	}
}

func printNginxVersions() {
	var versions []string
	versions = append(versions, versionsGenNginx()...)
	versions = append(versions, versionsGenOpenResty()...)
	versions = append(versions, versionsGenFreenginx()...)
	for _, v := range versions {
		fmt.Println(v)
	}
}

func versionCheck(version string) {
	if len(version) == 0 {
		log.Println("[warn]nginx version is not set.")
		log.Printf("[warn]nginx-build use %s.\n", builder.NginxVersion)
	}
}
