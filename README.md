# nginx-build

`nginx-build` - provides a command to build nginx seamlessly.

![gif](https://raw.githubusercontent.com/cubicdaiya/nginx-build/master/images/nginx-build.gif)

## Requirements

 * [git](https://git-scm.com/) and [hg](https://www.mercurial-scm.org/) for downloading 3rd party modules
 * [patch](https://savannah.gnu.org/projects/patch/) for applying patch to nginx

## Build Support

 * [nginx](https://nginx.org/)
 * [OpenResty](https://openresty.org/)
 * [freenginx](https://freenginx.org/)

## Installation

```bash
go install github.com/cubicdaiya/nginx-build@latest
```

## Quick Start

```console
nginx-build -d work
```

## Custom Configuration

`nginx-build` provides a mechanism for customizing configuration for building nginx.

### Configuration for building nginx

Prepare a configure script like the following.

```bash
#!/bin/sh

./configure \
--sbin-path=/usr/sbin/nginx \
--conf-path=/etc/nginx/nginx.conf \
```

Give this file to `nginx-build` with `-c`.

```bash
$ nginx-build -d work -c configure.example
```

### Embedding zlib statically

Give `-zlib` to `nginx-build`.

```bash
$ nginx-build -d work -zlib
```

`-zlibversion` is an option to set a version of zlib.

### Embedding PCRE statically

Give `-pcre` to `nginx-build`.

```bash
$ nginx-build -d work -pcre
```

`-pcreversion` is an option to set a version of PCRE.

### Embedding OpenSSL statically

Give `-openssl` to `nginx-build`.

```bash
$ nginx-build -d work -openssl
```

`-opensslversion` is an option to set a version of OpenSSL.

### Embedding LibreSSL statically

Give `-libressl` to `nginx-build`.

```bash
$ nginx-build -d work -libressl
```

`-libresslversion` is an option to set a version of LibreSSL.

### Using Custom SSL Libraries (e.g., BoringSSL)

`nginx-build` supports using arbitrary SSL libraries through custom SSL options. This is useful for libraries like BoringSSL that are not available as standard options.

#### Basic Usage

To use a custom SSL library, provide the download URL with `-customssl`. The URL can be either a tarball or a Git repository:

```bash
# Using BoringSSL from Git repository
$ nginx-build -d work -customssl https://boringssl.googlesource.com/boringssl -customsslname boringssl

# Using a tarball URL
$ nginx-build -d work -customssl https://example.com/customssl-1.0.0.tar.gz -customsslname customssl
```

#### Using a Specific Git Tag or Branch

For Git repositories, you can specify a tag or branch with `-customssltag`:

```bash
# Use BoringSSL with chromium-stable branch
$ nginx-build -d work \
  -customssl https://boringssl.googlesource.com/boringssl \
  -customsslname boringssl \
  -customssltag chromium-stable

# Use OpenSSL from Git with specific tag
$ nginx-build -d work \
  -customssl https://github.com/openssl/openssl.git \
  -customsslname openssl-git \
  -customssltag openssl-3.5.1

# Use oqs-provider (OpenSSL provider) with tag 0.9.0
$ nginx-build -d work \
  -customssl https://github.com/open-quantum-safe/oqs-provider.git \
  -customsslname oqs-provider \
  -customssltag 0.9.0
```

#### Using a Tarball URL

You can also use tarball URLs for custom SSL libraries:

```bash
# Using a custom OpenSSL build from a tarball
$ nginx-build -d work \
  -customssl https://example.com/myssl-1.0.0.tar.gz \
  -customsslname myssl
```

#### Available Options

- `-customssl`: URL of the custom SSL library (supports both Git repositories and tarballs)
- `-customsslname`: Name for the custom SSL library (used in directory names)
- `-customssltag`: Git tag or branch to checkout (only for Git repositories)

#### Supported Git Repository Formats

The following Git repository URL formats are automatically detected:
- URLs ending with `.git`
- URLs using `git://` protocol
- GitHub repository URLs (e.g., `https://github.com/user/repo`)
- Google Source URLs (e.g., `https://boringssl.googlesource.com/boringssl`)

Note: URLs containing `/releases/download/` or `/archive/` are treated as tarball downloads, not Git repositories.

### Embedding 3rd-party modules

`nginx-build` provides a mechanism for embedding 3rd-party modules.
Prepare a json file below.

```ini
[
  {
    "name": "ngx_http_hello_world",
    "form": "git",
    "url": "https://github.com/cubicdaiya/ngx_http_hello_world"
  }
]
```

Give this file to `nginx-build` with `-m`.

```bash
$ nginx-build -d work -m modules.json.example
```

#### Embedding 3rd-party module dynamically

Give `true` to `dynamic`.

```ini
[
  {
    "name": "ngx_http_hello_world",
    "form": "git",
    "url": "https://github.com/cubicdaiya/ngx_http_hello_world",
    "dynamic": true
  }
]
```

#### Provision for 3rd-party module

There are some 3rd-party modules expected provision. `nginx-build` provides the options such as `shprov` and `shprovdir` for this problem.
There is the example configuration below.

```ini
[
  {
    "name": "njs/nginx",
    "form": "hg",
    "url": "https://hg.nginx.org/njs",
    "shprov": "./configure && make",
    "shprovdir": ".."
  }
]
```

## Applying patch before building nginx

`nginx-build` provides the options such as `-patch` and `-patch-opt` for applying patches to the extracted sources. By default `-patch <path>` targets the primary tree (nginx, OpenResty, or freenginx depending on the build). To patch other bundled components prefix the flag with the component name, for example `-patch openssl=/path/to/openssl.patch`. Supported targets include `nginx`, `openresty`, `freenginx`, `pcre`/`pcre2`, `openssl`, `libressl`, `customssl` (or the custom name), and `zlib`.

```console
nginx-build \
 -d work \
 -openssl \
 -patch patches/nginx.patch \
 -patch openssl=patches/openssl.patch \
 -patch-opt "-p1"
```

## Idempotent build

`nginx-build` supports a certain level of idempotent build of nginx.
If you want to ensure a build of nginx idempotent and do not want to build nginx as same as already installed nginx,
give `-idempotent` to `nginx-build`.

```bash
$ nginx-build -d work -idempotent
```

`-idempotent` ensures an idempotent by checking the software versions below.

* nginx
* PCRE
* zlib
* OpenSSL

On the other hand, `-idempotent` does not cover versions of 3rd party modules and dynamic linked libraries.

## Build OpenResty

`nginx-build` supports to build [OpenResty](https://openresty.org/).

```bash
$ nginx-build -d work -openresty -pcre -openssl
```

If you don't install PCRE and OpenSSL on your system, it is required to add the option `-pcre` and `-openssl`.


And there is the limitation for the support of OpenResty.
`nginx-build` does not allow to use OpenResty's unique configure options directly.
If you want to use OpenResty's unique configure option, [Configuration for building nginx](#configuration-for-building-nginx) is helpful.

## Build freenginx

`nginx-build` supports to build [freenginx](https://freenginx.org/).

```bash
$ nginx-build -d work -freenginx -openssl
```

If you don't install OpenSSL on your system, it is required to add the option `-openssl`.
