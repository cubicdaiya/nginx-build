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
	versionsMinor0 := []int{
		45, // 0.1.x
		6,  // 0.2.x
		61, // 0.3.x
		14, // 0.4.x
		38, // 0.5.x
		39, // 0.6.x
		69, // 0.7.x
		55, // 0.8.x
		7,  // 0.9.x
	}
	// 1.x.x
	versionsMinor1 := []int{
		15, // 1.0.x
		19, // 1.1.x
		9,  // 1.2.x
		16, // 1.3.x
		7,  // 1.4.x
		13, // 1.5.x
		3,  // 1.6.x
		12, // 1.7.x
		1,  // 1.8.x
		15, // 1.9.x
		3,  // 1.10.x
		13, // 1.11.x
		2,  // 1.12.x
		12, // 1.13.x
		2,  // 1.14.x
		12, // 1.15.x
		1,  // 1.16.x
		9,  // 1.17.x
		0,  // 1.18.x
		2,  // 1.19.x
	}

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
	return []string{
		fmt.Sprintf("openresty-%s", builder.OpenRestyVersion),
	}
}

func versionsGenTengine() []string {
	return []string{
		fmt.Sprintf("tengine-%s", builder.TengineVersion),
	}
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
		log.Printf("[warn]nginx-build use %s.\n", builder.NginxVersion)
	}
}
