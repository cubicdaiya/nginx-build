package builder

import (
	"fmt"
	"testing"
)

func setupBuilders(t *testing.T) []Builder {
	builders := make([]Builder, ComponentMax)
	builders[ComponentNginx] = MakeBuilder(ComponentNginx, NginxVersion)
	builders[ComponentPcre] = MakeLibraryBuilder(ComponentPcre, PcreVersion, false)
	builders[ComponentOpenSSL] = MakeLibraryBuilder(ComponentOpenSSL, OpenSSLVersion, true)
	builders[ComponentLibreSSL] = MakeLibraryBuilder(ComponentLibreSSL, LibreSSLVersion, true)
	builders[ComponentZlib] = MakeLibraryBuilder(ComponentZlib, ZlibVersion, false)
	builders[ComponentOpenResty] = MakeBuilder(ComponentOpenResty, OpenRestyVersion)
	builders[ComponentTengine] = MakeBuilder(ComponentTengine, TengineVersion)
	return builders
}

func TestName(t *testing.T) {
	builders := setupBuilders(t)

	tests := []struct {
		got  string
		want string
	}{
		{
			got:  builders[ComponentNginx].name(),
			want: "nginx",
		},
		{
			got:  builders[ComponentPcre].name(),
			want: "pcre2",
		},
		{
			got:  builders[ComponentOpenSSL].name(),
			want: "openssl",
		},
		{
			got:  builders[ComponentLibreSSL].name(),
			want: "libressl",
		},
		{
			got:  builders[ComponentZlib].name(),
			want: "zlib",
		},
		{
			got:  builders[ComponentOpenResty].name(),
			want: "openresty",
		},
		{
			got:  builders[ComponentTengine].name(),
			want: "tengine",
		},
	}

	for _, test := range tests {
		if test.got != test.want {
			t.Fatalf("got: %v, want: %v", test.got, test.want)
		}
	}
}

func TestOption(t *testing.T) {
	builders := setupBuilders(t)

	tests := []struct {
		got  string
		want string
	}{
		{
			got:  builders[ComponentPcre].option(),
			want: "--with-pcre",
		},
		{
			got:  builders[ComponentOpenSSL].option(),
			want: "--with-openssl",
		},
		{
			got:  builders[ComponentLibreSSL].option(),
			want: "--with-openssl",
		},
		{
			got:  builders[ComponentZlib].option(),
			want: "--with-zlib",
		},
	}

	for _, test := range tests {
		if test.got != test.want {
			t.Fatalf("got: %v, want: %v", test.got, test.want)
		}
	}
}

func TestDownloadURL(t *testing.T) {
	builders := setupBuilders(t)

	tests := []struct {
		got  string
		want string
	}{
		{
			got:  builders[ComponentNginx].DownloadURL(),
			want: fmt.Sprintf("%s/nginx-%s.tar.gz", NginxDownloadURLPrefix, NginxVersion),
		},
		{
			got:  builders[ComponentPcre].DownloadURL(),
			want: fmt.Sprintf("%s/pcre2-%s/pcre2-%s.tar.gz", PcreDownloadURLPrefix, PcreVersion, PcreVersion),
		},
		{
			got:  builders[ComponentOpenSSL].DownloadURL(),
			want: fmt.Sprintf("%s/openssl-%s.tar.gz", OpenSSLDownloadURLPrefix, OpenSSLVersion),
		},
		{
			got:  builders[ComponentLibreSSL].DownloadURL(),
			want: fmt.Sprintf("%s/libressl-%s.tar.gz", LibreSSLDownloadURLPrefix, LibreSSLVersion),
		},
		{
			got:  builders[ComponentZlib].DownloadURL(),
			want: fmt.Sprintf("%s/zlib-%s.tar.gz", ZlibDownloadURLPrefix, ZlibVersion),
		},
		{
			got:  builders[ComponentOpenResty].DownloadURL(),
			want: fmt.Sprintf("%s/openresty-%s.tar.gz", OpenRestyDownloadURLPrefix, OpenRestyVersion),
		},
		{
			got:  builders[ComponentTengine].DownloadURL(),
			want: fmt.Sprintf("%s/tengine-%s.tar.gz", TengineDownloadURLPrefix, TengineVersion),
		},
	}

	for _, test := range tests {
		if test.got != test.want {
			t.Fatalf("got: %v, want: %v", test.got, test.want)
		}
	}
}

func TestSourcePath(t *testing.T) {
	builders := setupBuilders(t)

	tests := []struct {
		got  string
		want string
	}{
		{
			got:  builders[ComponentNginx].SourcePath(),
			want: fmt.Sprintf("nginx-%s", NginxVersion),
		},
		{
			got:  builders[ComponentPcre].SourcePath(),
			want: fmt.Sprintf("pcre2-%s", PcreVersion),
		},
		{
			got:  builders[ComponentOpenSSL].SourcePath(),
			want: fmt.Sprintf("openssl-%s", OpenSSLVersion),
		},
		{
			got:  builders[ComponentLibreSSL].SourcePath(),
			want: fmt.Sprintf("libressl-%s", LibreSSLVersion),
		},
		{
			got:  builders[ComponentZlib].SourcePath(),
			want: fmt.Sprintf("zlib-%s", ZlibVersion),
		},
		{
			got:  builders[ComponentOpenResty].SourcePath(),
			want: fmt.Sprintf("openresty-%s", OpenRestyVersion),
		},
		{
			got:  builders[ComponentTengine].SourcePath(),
			want: fmt.Sprintf("tengine-%s", TengineVersion),
		},
	}

	for _, test := range tests {
		if test.got != test.want {
			t.Fatalf("got: %v, want: %v", test.got, test.want)
		}
	}
}

func TestArchivePath(t *testing.T) {
	builders := setupBuilders(t)

	tests := []struct {
		got  string
		want string
	}{
		{
			got:  builders[ComponentNginx].ArchivePath(),
			want: fmt.Sprintf("nginx-%s.tar.gz", NginxVersion),
		},
		{
			got:  builders[ComponentPcre].ArchivePath(),
			want: fmt.Sprintf("pcre2-%s.tar.gz", PcreVersion),
		},
		{
			got:  builders[ComponentOpenSSL].ArchivePath(),
			want: fmt.Sprintf("openssl-%s.tar.gz", OpenSSLVersion),
		},
		{
			got:  builders[ComponentLibreSSL].ArchivePath(),
			want: fmt.Sprintf("libressl-%s.tar.gz", LibreSSLVersion),
		},
		{
			got:  builders[ComponentZlib].ArchivePath(),
			want: fmt.Sprintf("zlib-%s.tar.gz", ZlibVersion),
		},
		{
			got:  builders[ComponentOpenResty].ArchivePath(),
			want: fmt.Sprintf("openresty-%s.tar.gz", OpenRestyVersion),
		},
		{
			got:  builders[ComponentTengine].ArchivePath(),
			want: fmt.Sprintf("tengine-%s.tar.gz", TengineVersion),
		},
	}

	for _, test := range tests {
		if test.got != test.want {
			t.Fatalf("got: %v, want: %v", test.got, test.want)
		}
	}
}

func TestLogPath(t *testing.T) {
	builders := setupBuilders(t)

	tests := []struct {
		got  string
		want string
	}{
		{
			got:  builders[ComponentNginx].LogPath(),
			want: fmt.Sprintf("nginx-%s.log", NginxVersion),
		},
		{
			got:  builders[ComponentPcre].LogPath(),
			want: fmt.Sprintf("pcre2-%s.log", PcreVersion),
		},
		{
			got:  builders[ComponentOpenSSL].LogPath(),
			want: fmt.Sprintf("openssl-%s.log", OpenSSLVersion),
		},
		{
			got:  builders[ComponentLibreSSL].LogPath(),
			want: fmt.Sprintf("libressl-%s.log", LibreSSLVersion),
		},
		{
			got:  builders[ComponentZlib].LogPath(),
			want: fmt.Sprintf("zlib-%s.log", ZlibVersion),
		},
		{
			got:  builders[ComponentOpenResty].LogPath(),
			want: fmt.Sprintf("openresty-%s.log", OpenRestyVersion),
		},
		{
			got:  builders[ComponentTengine].LogPath(),
			want: fmt.Sprintf("tengine-%s.log", TengineVersion),
		},
	}

	for _, test := range tests {
		if test.got != test.want {
			t.Fatalf("got: %v, want: %v", test.got, test.want)
		}
	}
}

func TestLibrary(t *testing.T) {
	builders := setupBuilders(t)

	tests := []struct {
		got  bool
		want bool
	}{
		{
			got:  builders[ComponentPcre].Static,
			want: false,
		},
		{
			got:  builders[ComponentOpenSSL].Static,
			want: true,
		},
		{
			got:  builders[ComponentLibreSSL].Static,
			want: true,
		},
		{
			got:  builders[ComponentZlib].Static,
			want: false,
		},
	}

	for _, test := range tests {
		if test.got != test.want {
			t.Fatalf("got: %v, want: %v", test.got, test.want)
		}
	}
}

func TestMakeStaticLibrary(t *testing.T) {
	builders := setupBuilders(t)

	tests := []struct {
		builder       Builder
		staticLibrary StaticLibrary
		version       string
	}{
		{
			builder:       builders[ComponentPcre],
			staticLibrary: MakeStaticLibrary(&builders[ComponentPcre]),
			version:       PcreVersion,
		},
		{
			builder:       builders[ComponentOpenSSL],
			staticLibrary: MakeStaticLibrary(&builders[ComponentOpenSSL]),
			version:       OpenSSLVersion,
		},
		{
			builder:       builders[ComponentLibreSSL],
			staticLibrary: MakeStaticLibrary(&builders[ComponentLibreSSL]),
			version:       LibreSSLVersion,
		},
		{
			builder:       builders[ComponentZlib],
			staticLibrary: MakeStaticLibrary(&builders[ComponentZlib]),
			version:       ZlibVersion,
		},
	}

	for _, test := range tests {
		if test.builder.name() != test.staticLibrary.Name {
			t.Fatalf("not equal name() between builder(%v) and static library(%v)", test.builder.name(), test.staticLibrary.Name)
		}
		if test.builder.option() != test.staticLibrary.Option {
			t.Fatalf("not equal option() between builder(%v) and static library(%v)", test.builder.option(), test.staticLibrary.Option)
		}
		if test.builder.Version != test.version {
			t.Fatalf("not equal version between builder's and default's")
		}
	}
}
