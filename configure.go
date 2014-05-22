package main

import (
        "fmt"
        "io/ioutil"
        "strings"
)

func configureModule3rdGen(modules3rd []Module3rd) string {
        result := ""
        for i := 0; i < len(modules3rd); i++ {
                result += fmt.Sprintf("--add-module=../%s \\\n", modules3rd[i].Name)
        }
        return result
}

func (builder *Builder) ConfigureGen(conf string, modules3rd []Module3rd, pcreStatic bool, pcreVersion string, openSSLStaic bool, openSSLVersion string) error {
        configure := `#!/bin/sh

./configure `

        if conf != "" {
                configure += "\\\n"
                options := strings.Split(conf, "\n")

                for i := 0; i < len(options); i++ {
                        options[i] += " \\"
                }

                conf = strings.Join(options, "\n")
                configure += conf
        }

        if pcreStatic {
                configure += fmt.Sprintf("--with-pcre=../pcre-%s \\\n", pcreVersion)
        }

        if openSSLStaic {
                configure += fmt.Sprintf("--with-openssl=../openssl-%s \\\n", openSSLVersion)
        }

        configure_modules3rd := configureModule3rdGen(modules3rd)
        configure += configure_modules3rd

        return ioutil.WriteFile("./nginx-configure", []byte(configure), 0644)
}
