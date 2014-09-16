package main

import (
	"fmt"
	"runtime"
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

	versionsMinor0 := []int{45, 6, 61, 14, 38, 39, 69, 55, 7} // 0.x.x
	versionsMinor1 := []int{15, 19, 9, 16, 7, 13, 0, 5}       // 1.x.x

	// 0.1.0 ~ 0.9.7
	for i := 0; i < 9; i++ {
		versions = append(versions, versionsSubmajorGen(0, i+1, versionsMinor0[i])...)
	}

	// 1.0.0 ~
	for i := 0; i < 8; i++ {
		versions = append(versions, versionsSubmajorGen(1, i, versionsMinor1[i])...)
	}

	return versions
}

func printNginxVersions() {
	versions := versionsGen()
	for _, v := range versions {
		fmt.Println(v)
	}
}

func printNginxBuildVersion() {
	fmt.Printf(`nginx-build %s
Compiler: %s %s
Copyright (C) 2014 Tatsuhiko Kubo <cubicdaiya@gmail.com>
`,
		NGINX_BUILD_VERSION,
		runtime.Compiler,
		runtime.Version())

}
