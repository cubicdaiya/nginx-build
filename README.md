# nginx-build

nginx-build - provides a command to build nginx seamlessly.

# Dependencies

 * [Go](http://golang.org/)
 * [github.com/robfig/config](https://github.com/robfig/config)

## Install

```
go get github.com/robfig/config
go get github.com/cubicdaiya/nginx-build
```

## Usage

```
mkdir -p ~/opt/nginx
nginx-build -v 1.7.0 -d ~/opt/nginx
cd ~/opt/nginx/1.7.0/nginx-1.7.0
objs/bin/nginx -V
```
