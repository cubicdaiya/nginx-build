# nginx-build

`nginx-build` - provides a command to build nginx seamlessly.

![gif](https://raw.githubusercontent.com/cubicdaiya/nginx-build/master/images/nginx-build.gif)

## Requirements

 * [git](https://git-scm.com/) and [hg](https://www.mercurial-scm.org/) for downloading 3rd party modules
 * [patch](https://savannah.gnu.org/projects/patch/) for applying patch to nginx

## Build Support

 * [nginx](https://nginx.org/)
 * [OpenResty](https://openresty.org/)
 * [Tengine](https://tengine.taobao.org/)

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

#### About `--add-module` and `--add-dynamic-module`

`nginx-build` allows to use multiple `--add-module` and `--add-dynamic-module`.

```bash
$ nginx-build \
-d work \
--add-module=/path/to/ngx_dynamic_upstream \
--add-dynamic-module=/path/to/ngx_small_light
```

On the other hand, `nginx-build` allows to embed multiple 3rd party modules with the single `--add-module` and `--add-dynamic-module` like the following, too.

```bash
$ nginx-build \
-d work \
--add-module=/path/to/ngx_small_light,/path/to/ngx_dynamic_upstream
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

### Embedding 3rd-party modules

`nginx-build` provides a mechanism for embedding 3rd-party modules.
Prepare a ini-file below.

```ini
[ngx_http_hello_world]
form=git
url=https://github.com/cubicdaiya/ngx_http_hello_world
```

Give this file to `nginx-build` with `-m`.

```bash
$ nginx-build -d work -m modules.cfg.example
```

#### Embedding 3rd-party module dynamically

Give `dynamic=true`.

```ini
[ngx_http_hello_world]
form=git
url=https://github.com/cubicdaiya/ngx_http_hello_world
dynamic=true
```

#### Provision for 3rd-party module

There are some 3rd-party modules expected provision. `nginx-build` provides the options such as `shprov` and `shprovdir` for this problem.
There is the example configuration below.

```ini
[njs/nginx]
form=hg
url=https://hg.nginx.org/njs
shprov=./configure && make
shprovdir=..
```

## Applying patch before building nginx

`nginx-build` provides the options such as `-patch` and `-patch-opt` for applying patch to nginx.

```console
nginx-build \
 -d work \
 -patch something.patch \
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
But you can use the common options of nginx and OpenResty directly.
If you want to use OpenResty's unique configure option, [Configuration for building nginx](#configuration-for-building-nginx) is helpful.

## Build Tengine

`nginx-build` supports to build [Tengine](https://tengine.taobao.org/).

```bash
$ nginx-build -d work -tengine -openssl
```

If you don't install OpenSSL on your system, it is required to add the option `-openssl`.

There is the limitation for the support of [Tengine](https://tengine.taobao.org/).
`nginx-build` does not allow to use Tengine's unique configure options directly.
But you can use the common options of nginx and Tengine directly.
