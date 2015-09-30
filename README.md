# nginx-build

[![Build Status](https://drone.io/github.com/cubicdaiya/nginx-build/status.png)](https://drone.io/github.com/cubicdaiya/nginx-build/latest)

`nginx-build` - provides a command to build nginx seamlessly.

## Requirements

 * [wget](https://www.gnu.org/software/wget/) for downloading nginx and external libraries
 * [git](https://git-scm.com/) for downloading 3rd party modules

## Build Support

 * [nginx](http://nginx.org/)
 * [OpenResty](http://www.openresty.com/)
 * [Tengine](http://tengine.taobao.org/)

## Installation

```bash
go get -u github.com/cubicdaiya/nginx-build
```

If you don't have go-runtime, you may download the binary from [here](https://github.com/cubicdaiya/nginx-build/releases).


If you are Mac OS X user, you can use [Homebrew](http://brew.sh/).

```
brew tap cubicdaiya/nginx-build
brew install nginx-build
```

## Quick Start

```console
$ nginx-build -d work
nginx-build: 0.4.4
Compiler: gc go1.5
2015/08/30 11:36:26 Download nginx-1.9.4.....
2015/08/30 11:36:28 Extract nginx-1.9.4.tar.gz.....
2015/08/30 11:36:28 Generate configure script for nginx-1.9.4.....
2015/08/30 11:36:28 Configure nginx-1.9.4.....
2015/08/30 11:36:35 Build nginx-1.9.4.....
2015/08/30 11:36:40 Complete building nginx!

nginx version: nginx/1.9.4
built by clang 6.1.0 (clang-602.0.53) (based on LLVM 3.6.0svn)
configure arguments: --with-cc-opt=-Wno-deprecated-declarations

2015/08/30 11:36:40 Enter the following command for install nginx.

   $ cd work/1.9.4/nginx-1.9.4
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

From `v0.4.0`, `nginx-build` allows to use nginx's configure options directly.

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
--with-http_spdy_module \
--with-http_ssl_module \
--with-pcre-jit \
```

#### About `--add-module`

`nginx-build` allows to use multiple `--add-module` from `v0.4.2`.

```bash
$ nginx-build \
-d work \
--add-module=../nginx/ngx_small_light \
--add-module=../nginx/ngx_dynamic_upstream
```

On the other hand, `nginx-build` allows to embed multiple 3rd party modules with the single `--add-module` like the following, too.

```bash
$ nginx-build \
-d work \
--add-module=../nginx/ngx_small_light,../nginx/ngx_dynamic_upstream
```

The values of `--add-module` of `nginx-build` must be separated by ','. This feature is introduced in `v0.4.0`.

#### Limitations

There are the limitations for the direct configuration below.

 * `--with-pcre`(force PCRE library usage) is not allowed
  * `--with-pcre=DIR`(set path to PCRE library sources) is allowed
 * `--with-libatomic`(force libatomic_ops library usage) is not allowed
  * `--with-libatomic=DIR`(set path to libatomic_ops library sources) is allowed

The limitations above are attributed by the flag package of Go. (multiple and different types from each other are not allowed)

### Embedding zlib statically

Give `-zlib` to `nginx-build`.

```bash
$ nginx-build -d work -zlib
```

`-zlibversion` is an option to set a version of zlib.

```bash
$ nginx-build -d work -zlib -zlibversion=1.2.8
```

### Embedding PCRE statically

Give `-pcre` to `nginx-build`.

```bash
$ nginx-build -d work -pcre
```

`-pcreversion` is an option to set a version of PCRE.

```bash
$ nginx-build -d work -pcre -pcreversion=8.37
```

### Embedding OpenSSL statically

Give `-openssl` to `nginx-build`.

```bash
$ nginx-build -d work -openssl
```

`-opensslversion` is an option to set a version of OpenSSL.

```bash
$ nginx-build -d work -openssl -opensslversion=1.0.2d
```

### Embedding 3rd-party modules

`nginx-build` provides a mechanism for embedding 3rd-party modules.
Prepare a ini-file below.

```ini
[echo-nginx-module]
form=git
url=https://github.com/openresty/echo-nginx-module.git
rev=v0.58

[ngx_devel_kit]
form=git
url=https://github.com/simpl/ngx_devel_kit
rev=v0.2.19
```

Give this file to `nginx-build` with `-m`.

```bash
$ nginx-build -d work -m modules.cfg.example
```

## Build OpenResty

`nginx-build` supports to build [OpenResty](http://www.openresty.com/).

```bash
$ nginx-build -d work -openresty -pcre
```

If you don't install PCRE on your system, it is required to add the option `-pcre`.


And there is the limitation for the support of OpenResty.
`nginx-build` does not allow to use OpenResty's unique configure options directly.
But you can use the common options of nginx and OpenResty directly.

## Build Tengine

`nginx-build` supports to build [Tengine](http://tengine.taobao.org/).

```bash
$ nginx-build -d work -tengine
```

There is the limitation for the support of [Tengine](http://tengine.taobao.org/).
`nginx-build` does not allow to use Tengine's unique configure options directly.
But you can use the common options of nginx and Tengine directly.
