package module3rd

import (
	"fmt"

	"github.com/cubicdaiya/nginx-build/util"
	"github.com/go-ini/ini"
)

func Load(path string) ([]Module3rd, error) {
	var modules []Module3rd
	if len(path) > 0 {
		if !util.FileExists(path) {
			return modules, fmt.Errorf("modulesConfPath(%s) does not exist.", path)
		}
		modulesConf, err := ini.Load(path)
		if err != nil {
			return modules, err
		}
		modules = loadModules(modulesConf)
	}
	return modules, nil
}

func loadModule(s *ini.Section) Module3rd {
	var (
		module Module3rd
	)

	module.Name = s.Name()
	module.Form = s.Key("form").String()
	if module.Form == "" {
		module.Form = "git"
	}
	module.Url = s.Key("url").String()
	module.Rev = s.Key("rev").String()
	module.Shprov = s.Key("shprov").String()
	module.ShprovDir = s.Key("shprovdir").String()
	module.Dynamic = s.Key("dynamic").MustBool()

	return module
}

func loadModules(f *ini.File) []Module3rd {
	var modules []Module3rd
	for _, s := range f.Sections() {
		if s.Name() == "DEFAULT" {
			continue
		}
		modules = append(modules, loadModule(s))
	}
	return modules
}
