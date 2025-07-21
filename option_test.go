package main

import (
	"testing"
)

func TestMakeNginxBuildOptions(t *testing.T) {
	options := makeNginxBuildOptions()

	// Test custom SSL options exist
	customSSLTests := []struct {
		name string
		key  string
		desc string
	}{
		{
			name: "customssl option",
			key:  "customssl",
			desc: "download URL for custom SSL library",
		},
		{
			name: "customsslname option",
			key:  "customsslname",
			desc: "name for custom SSL library",
		},
		{
			name: "customssltag option",
			key:  "customssltag",
			desc: "git tag/branch for custom SSL library",
		},
	}

	for _, test := range customSSLTests {
		t.Run(test.name, func(t *testing.T) {
			if opt, ok := options.Values[test.key]; !ok {
				t.Errorf("Option %s not found in Values", test.key)
			} else if opt.Desc != test.desc {
				t.Errorf("Option %s description = %q, want %q", test.key, opt.Desc, test.desc)
			}
		})
	}
}

func TestIsNginxBuildOption(t *testing.T) {
	// Initialize global options
	nginxBuildOptions = makeNginxBuildOptions()

	tests := []struct {
		name string
		key  string
		want bool
	}{
		{
			name: "valid value option - customssl",
			key:  "customssl",
			want: true,
		},
		{
			name: "valid value option - customsslname",
			key:  "customsslname",
			want: true,
		},
		{
			name: "valid bool option - openssl",
			key:  "openssl",
			want: true,
		},
		{
			name: "invalid option",
			key:  "invalid-option",
			want: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := isNginxBuildOption(test.key); got != test.want {
				t.Errorf("isNginxBuildOption(%q) = %v, want %v", test.key, got, test.want)
			}
		})
	}
}
