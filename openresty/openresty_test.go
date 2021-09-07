package openresty

import (
	"testing"
)

func TestOpenrestyName(t *testing.T) {

	tests := []struct {
		version string
		name    string
	}{
		{
			version: "1.9.7.2",
			name:    "ngx_openresty",
		},
		{
			version: "1.9.7.3",
			name:    "openresty",
		},
		{
			version: "1.9.7.4",
			name:    "openresty",
		},
		{
			version: "1.15.8.1rc1",
			name:    "openresty",
		},
	}

	for _, test := range tests {
		name := Name(test.version)
		if name != test.name {
			t.Fatalf("got: %v, want: %v", name, test.name)
		}
	}
}
