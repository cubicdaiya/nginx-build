package main

import (
	"testing"
)

func TestIsGitURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want bool
	}{
		{
			name: "Git URL with .git suffix",
			url:  "https://github.com/user/repo.git",
			want: true,
		},
		{
			name: "Git URL with git:// protocol",
			url:  "git://github.com/user/repo",
			want: true,
		},
		{
			name: "GitHub URL without .git",
			url:  "https://github.com/user/repo",
			want: true,
		},
		{
			name: "GitHub releases URL",
			url:  "https://github.com/openssl/openssl/releases/download/openssl-3.5.1/openssl-3.5.1.tar.gz",
			want: false,
		},
		{
			name: "GitHub archive URL",
			url:  "https://github.com/user/repo/archive/v1.0.tar.gz",
			want: false,
		},
		{
			name: "Google Source URL",
			url:  "https://boringssl.googlesource.com/boringssl",
			want: true,
		},
		{
			name: "Regular tar.gz URL",
			url:  "https://example.com/library.tar.gz",
			want: false,
		},
		{
			name: "Regular tar.bz2 URL",
			url:  "https://example.com/library.tar.bz2",
			want: false,
		},
		{
			name: "Empty URL",
			url:  "",
			want: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := isGitURL(test.url); got != test.want {
				t.Errorf("isGitURL(%q) = %v, want %v", test.url, got, test.want)
			}
		})
	}
}
