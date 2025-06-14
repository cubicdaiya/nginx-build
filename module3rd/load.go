package module3rd

import (
	"encoding/json"
	"fmt"
	"os"
)

func Load(path string) ([]Module3rd, error) {
	var modules []Module3rd
	if len(path) > 0 {
		f, err := os.Open(path)
		if err != nil {
			return modules, err
		}
		defer f.Close()
		if err := json.NewDecoder(f).Decode(&modules); err != nil {
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
