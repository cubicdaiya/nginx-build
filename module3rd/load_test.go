package module3rd

import (
	"testing"
)

func TestModules3rd(t *testing.T) {

	modules3rdConf := "../config/modules.cfg.example"
	modules3rd, err := Load(modules3rdConf)
	if err != nil {
		t.Fatalf("Failed to load %s", modules3rdConf)
	}

	for _, m := range modules3rd {
		var want Module3rd
		switch m.Name {
		case "headers-more-nginx-module":
			want.Name = "headers-more-nginx-module"
			want.Form = "git"
			want.Url = "https://github.com/openresty/headers-more-nginx-module.git"
			want.Rev = "v0.32"
			want.Dynamic = false
		case "ngx_devel_kit":
			want.Name = "ngx_devel_kit"
			want.Form = "git"
			want.Url = "https://github.com/simpl/ngx_devel_kit"
			want.Rev = "v0.3.0"
			want.Dynamic = false
		case "ngx_small_light":
			want.Name = "ngx_small_light"
			want.Form = "git"
			want.Url = "https://github.com/cubicdaiya/ngx_small_light"
			want.Rev = "v0.9.2"
			want.Shprov = "./setup --with-gd"
			want.Dynamic = true
		default:
			t.Fatalf("unexpected module: %v", m)
		}

		if m != want {
			t.Fatalf("got: %v, want: %v", m, want)
		}
	}
}
