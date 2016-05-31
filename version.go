package main

import (
	"fmt"
	"log"

	"github.com/cubicdaiya/nginx-build/builder"
)

func versionsSubmajorGen(major, submajor, minor int) []string {
	var versions []string

	for i := 0; i <= minor; i++ {
		v := fmt.Sprintf("%d.%d.%d", major, submajor, i)
		versions = append(versions, v)
	}
	return versions
}

func versionsGen() []string {
	var versions []string
	// 0.x.x
	versionsMinor0 := []int{45, 6, 61, 14, 38, 39, 69, 55, 7}
	// 1.x.x
	versionsMinor1 := []int{15, 19, 9, 16, 7, 13, 3, 12, 1, 15, 1, 1}

	// 0.1.0 ~ 0.9.7
	for i := 0; i < len(versionsMinor0); i++ {
		versions = append(versions, versionsSubmajorGen(0, i+1, versionsMinor0[i])...)
	}

	// 1.0.0 ~
	for i := 0; i < len(versionsMinor1); i++ {
		versions = append(versions, versionsSubmajorGen(1, i, versionsMinor1[i])...)
	}

	return versions
}

func versionsGenOpenResty() []string {
	var versions []string

	versions = append(versions, fmt.Sprintf("openresty-%s", builder.OPENRESTY_VERSION))

	return versions
}

func versionsGenTengine() []string {
	var versions []string

	versions = append(versions, fmt.Sprintf("tengine-%s", builder.TENGINE_VERSION))

	return versions
}

func printNginxVersions() {
	versions := versionsGen()
	versions = append(versions, versionsGenOpenResty()...)
	versions = append(versions, versionsGenTengine()...)
	for _, v := range versions {
		fmt.Println(v)
	}
}

func versionCheck(version string) {
	if len(version) == 0 {
		log.Println("[warn]nginx version is not set.")
		log.Printf("[warn]nginx-build use %s.\n", builder.NGINX_VERSION)
	}
}
