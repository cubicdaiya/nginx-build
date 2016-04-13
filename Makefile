VERSION=0.8.0

nginx-build: *.go
	gom build -o nginx-build

build-cross:
	GOOS=linux GOARCH=amd64 gom build -ldflags '-s -w' -o bin/linux/amd64/nginx-build
	GOOS=darwin GOARCH=amd64 gom build -ldflags '-s -w' -o bin/darwin/amd64/nginx-build

dist: build-cross
	cd bin/linux/amd64/ && tar zcvf nginx-build-linux-amd64-${VERSION}.tar.gz nginx-build
	cd bin/darwin/amd64/ && tar zcvf nginx-build-darwin-amd64-${VERSION}.tar.gz nginx-build

# ImageMagick and GD are required for ngx_small_light
build-example: nginx-build
	./nginx-build -c config/configure.example -m config/modules.cfg.example -d work -clear

gom:
	go get -u github.com/mattn/gom

bundle:
	gom install

check:
	gom test

fmt:
	go fmt ./...

install:
	install nginx-build /usr/local/bin/nginx-build
	install doc/_build/man/nginx-build.7 /usr/local/share/man/man7/nginx-build.7

.PHONY: doc
doc:
	cd doc; make man

clean:
	rm -rf nginx-build
	rm -rf bin/linux/amd64/nginx-build*
	rm -rf bin/darwin/amd64/nginx-build*
