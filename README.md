# nginx-build

`nginx-build` - provides a command to build nginx seamlessly.

## Setup

```
go get -u github.com/robfig/config
go get -u github.com/cubicdaiya/nginx-build
```

## Quick Start

```
mkdir -p ~/opt/nginx
nginx-build -v 1.7.0 -d ~/opt/nginx
cd ~/opt/nginx/1.7.0/nginx-1.7.0
objs/bin/nginx -V
```

## Custom Configuration

### Configuration for building nginx

`nginx-build` provides a mechanism for custom configuration for building nginx.
Prepare a configure script like the following.

```
#!bin/sh

./configure \
--sbin-path=/usr/sbin/nginx \
--conf-path=/etc/nginx/nginx.conf \
--with-cc-opt="-Wno-deprecated-declarations" \
```

Give this file to `nginx-build` with `-c`.

```
$ nginx-build -v 1.7.0 -d ~/opt/nginx -c configure.example
```

#### Caution about `--with-zlib`

Don't use `--with-zlib` for embedding ZLIB statically.
Instead you should use `-zlib` and `-zlibversion`.

#### Caution about `--with-pcre`

Don't use `--with-pcre` for embedding PCRE statically.
Instead you should use `-pcre` and `-pcreversion`.

#### Caution about `--with-openssl`

Don't use `--with-openssl` for embedding OpenSSL in this text-file.
Instead you should use `-openssl` and `-opensslversion`.

#### Caution about `--add-module`

Don't use `--add-module` for embedding 3rd-party module in this text-file.
Instead you should use `-m`.

### Embedding ZLIB statically

Give `-zlib` to `nginx-build`.

```
$ nginx-build -v 1.7.0 -d ~/opt/nginx -zlib
```

`-zlibverson` is a option to set a version of ZLIB.

```
$ nginx-build -v 1.7.0 -d ~/opt/nginx -zlib -zlibversion=1.2.8
```

### Embedding PCRE statically

Give `-pcre` to `nginx-build`.

```
$ nginx-build -v 1.7.0 -d ~/opt/nginx -pcre
```

`-pcreverson` is a option to set a version of PCRE.

```
$ nginx-build -v 1.7.0 -d ~/opt/nginx -pcre -pcreversion=8.35
```

### Embedding OpenSSL statically

Give `-openssl` to `nginx-build`.

```
$ nginx-build -v 1.7.0 -d ~/opt/nginx -openssl
```

`-opensslverson` is a option to set a version of OPENSSL.

```
$ nginx-build -v 1.7.0 -d ~/opt/nginx -openssl -opensslversion=1.0.1g
```

### Embedding 3rd-party modules

`nginx-build` provides a mechanism for embedding 3rd-party modules.
Prepare a ini-file like the following.

```
[echo-nginx-module]
form=github
url=https://github.com/openresty/echo-nginx-module.git
rev=v0.53

[ngx_devel_kit]
form=github
url=https://github.com/simpl/ngx_devel_kit
rev=v0.2.19
```

Give this file to `nginx-build` with `-m`.

```
$ nginx-build -v 1.7.0 -d ~/opt/nginx -m modules.cfg.example
```
