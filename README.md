# nginx-build

nginx-build - provides a command to build nginx seamlessly.

## Setup

```
git clone git@github.com:cubicdaiya/nginx-build.git
cd nginx-build
gom install
gom build
go install
```

## Usage

```
mkdir -p ~/opt/nginx
nginx-build -v 1.7.0 -d ~/opt/nginx
cd ~/opt/nginx/1.7.0/nginx-1.7.0
objs/bin/nginx -V
```
