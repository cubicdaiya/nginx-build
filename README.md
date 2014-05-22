# nginx-build

`nginx-build` - provides a command to build nginx seamlessly.

## Setup

```
git clone git@github.com:cubicdaiya/nginx-build.git
cd nginx-build
gom install
gom build
go install
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
Prepare a text-file like the following.

```
--sbin-path=/usr/sbin/nginx
--conf-path=/etc/nginx/nginx.conf
--error-log-path=/var/log/nginx/error.log
--pid-path=/var/run/nginx.pid
--lock-path=/var/lock/nginx.lock
--http-log-path=/var/log/nginx/access.log
--http-client-body-temp-path=/var/lib/nginx/body
--http-proxy-temp-path=/var/lib/nginx/proxy
--with-http_stub_status_module
--http-fastcgi-temp-path=/var/lib/nginx/fastcgi
--with-debug
--with-pcre-jit
--with-http_spdy_module
--with-http_ssl_module
--with-cc-opt="-Wno-deprecated-declarations"
```

Give this file to `nginx-build` with `-c`.

```
$ nginx-build -v 1.7.0 -d ~/opt/nginx -c configure.options.example
```

#### Caution about `--with-pcre`

Don't use `--with-pcre` for embedding PCRE statically.
Instead you should use `-pcre` and `-pcreversion`.

#### Caution about `--add-module`

Don't use `--add-module` for embedding 3rd-party module in this text-file.
Instead you should use `-m`.

### Embedding PCRE statically

Give `-pcre` to `nginx-build`.

```
$ nginx-build -v 1.7.0 -d ~/opt/nginx -pcre
```

`-pcreverson` is a option to set a version of PCRE.

```
$ nginx-build -v 1.7.0 -d ~/opt/nginx -pcre -pcreversion=8.35
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
