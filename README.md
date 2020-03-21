# nginx-build

[![Build Status](https://travis-ci.org/cubicdaiya/nginx-build.svg?branch=master)](https://travis-ci.org/cubicdaiya/nginx-build)

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
go get -u github.com/cubicdaiya/nginx-build
```

If you don't have go-runtime, you may download the binary from [here](https://github.com/cubicdaiya/nginx-build/releases).


If you are Mac OS X user, you can use [Homebrew](https://brew.sh/).

```
brew tap cubicdaiya/nginx-build
brew install nginx-build
```

## Quick Start

```console
$ nginx-build -d work
nginx-build: 0.9.1
Compiler: gc go1.6.1
2016/04/21 09:40:12 Download nginx-1.9.15.....
2016/04/21 09:40:19 Extract nginx-1.9.15.tar.gz.....
2016/04/21 09:40:19 Generate configure script for nginx-1.9.15.....
2016/04/21 09:40:19 Configure nginx-1.9.15.....
2016/04/21 09:40:22 Build nginx-1.9.15.....
2016/04/21 09:40:25 Complete building nginx!

nginx version: nginx/1.9.15
built by gcc 4.8.4 (Ubuntu 4.8.4-2ubuntu1~14.04.1)
configure arguments:

2016/04/21 09:40:25 Enter the following command for install nginx.

   $ cd work/nginx/1.9.15/nginx-1.9.15
   $ sudo make install

$
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

### Direct configuration for building nginx

In the `v0.4.0` or later, `nginx-build` allows to use nginx's configure options directly.

```bash
$ nginx-build -d work \
--sbin-path=/usr/sbin/nginx \
--conf-path=/etc/nginx/nginx.conf \
--error-log-path=/var/log/nginx/error.log \
--pid-path=/var/run/nginx.pid \
--lock-path=/var/lock/nginx.lock \
--http-log-path=/var/log/nginx/access.log \
--http-client-body-temp-path=/var/lib/nginx/body \
--http-proxy-temp-path=/var/lib/nginx/proxy \
--with-http_stub_status_module \
--http-fastcgi-temp-path=/var/lib/nginx/fastcgi \
--with-debug \
--with-http_gzip_static_module \
--with-http_v2_module \
--with-http_ssl_module \
--with-pcre-jit \
```

But there are limitations. See [here](https://github.com/cubicdaiya/nginx-build#limitations) about details.

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

#### Limitations

There are the limitations for the direct configuration below.

 * `--with-pcre`(force PCRE library usage) is not allowed
  * `--with-pcre=DIR`(set path to PCRE library sources) is allowed
 * `--with-libatomic`(force libatomic_ops library usage) is not allowed
  * `--with-libatomic=DIR`(set path to libatomic_ops library sources) is allowed

The limitations above are attributed by the flag package of Go. (multiple and different types from each other are not allowed)
By the way, the options above are allowed in [a prepared configure script](https://github.com/cubicdaiya/nginx-build#configuration-for-building-nginx), of course.

### Embedding zlib statically

Give `-zlib` to `nginx-build`.

```bash
$ nginx-build -d work -zlib
```

`-zlibversion` is an option to set a version of zlib.

```bash
$ nginx-build -d work -zlib -zlibversion=1.2.9
```

### Embedding PCRE statically

Give `-pcre` to `nginx-build`.

```bash
$ nginx-build -d work -pcre
```

`-pcreversion` is an option to set a version of PCRE.

```bash
$ nginx-build -d work -pcre -pcreversion=8.38
```

### Embedding OpenSSL statically

Give `-openssl` to `nginx-build`.

```bash
$ nginx-build -d work -openssl
```

`-opensslversion` is an option to set a version of OpenSSL.

```bash
$ nginx-build -d work -openssl -opensslversion=1.0.2f
```

### Embedding LibreSSL statically

Give `-libressl` to `nginx-build`.

```bash
$ nginx-build -d work -libressl
```

`-libresslversion` is an option to set a version of LibreSSL.

```bash
$ nginx-build -d work -libressl -libresslversion=2.9.2
```

Known issue, the build with libressl fails on MacOSX.

### Embedding 3rd-party modules

`nginx-build` provides a mechanism for embedding 3rd-party modules.
Prepare a ini-file below.

```ini
[ngx_devel_kit]
form=git
url=https://github.com/simpl/ngx_devel_kit
rev=v0.2.19
```

Give this file to `nginx-build` with `-m`.

```bash
$ nginx-build -d work -m modules.cfg.example
```

#### Embedding 3rd-party module dynamically

Give `dynamic=true`.

```ini
[ngx_dynamic_upstream]
form=git
url=https://github.com/cubicdaiya/ngx_dynamic_upstream.git
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

[ngx_small_light]
form=git
url=https://github.com/cubicdaiya/ngx_small_light
rev=v0.9.2
dynamic=true
shprov=./setup
```

## Applying patch before building nginx

`nginx-build` provides the options such as `-patch` and `-patch-opt` for applying patch to nginx.

```console
nginx-build \
 -d work \
 -patch nginx__http2_spdy.patch \
 -patch-opt "-p1" \
 -v 1.9.7 \
 --with-http_spdy_module \
 --with-http_v2_module
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
