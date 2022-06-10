package module3rd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/cubicdaiya/nginx-build/util"
)

func Load(path string) ([]Module3rd, error) {
	var modules []Module3rd
	if len(path) > 0 {
		if !util.FileExists(path) {
			return modules, fmt.Errorf("modulesConfPath(%s) does not exist.", path)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return modules, err
		}
		if err := json.Unmarshal(data, &modules); err != nil {
			return modules, fmt.Errorf("modulesConfPath(%s) is invalid JSON.", path)
		}
		for i, _ := range modules {
			if modules[i].Form == "" {
				modules[i].Form = "git"
			}
		}
	}
	return modules, nil
}
