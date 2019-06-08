VERSION=0.11.6

nginx-build: *.go builder/*.go command/*.go configure/*.go module3rd/*.go openresty/*.go util/*.go
	GO111MODULE=on go build -ldflags '-X main.NginxBuildVersion=${VERSION}' -o $@

build-cross:
	GO111MODULE=on GOOS=linux GOARCH=amd64 go build -ldflags '-s -w -X main.NginxBuildVersion=${VERSION}' -o bin/linux/amd64/nginx-build
	GO111MODULE=on GOOS=darwin GOARCH=amd64 go build -ldflags '-s -w -X main.NginxBuildVersion=${VERSION}' -o bin/darwin/amd64/nginx-build

dist: build-cross
	cd bin/linux/amd64/ && tar cvf nginx-build-linux-amd64-${VERSION}.tar nginx-build && zopfli nginx-build-linux-amd64-${VERSION}.tar
	cd bin/darwin/amd64/ && tar cvf nginx-build-darwin-amd64-${VERSION}.tar nginx-build && zopfli nginx-build-darwin-amd64-${VERSION}.tar

.PHONY: release
release:
	mkdir release
	GO111MODULE=on GOOS=linux go build -ldflags '-s -w -X main.NginxBuildVersion=${VERSION}' -o nginx-build
	tar cvzf release/nginx-build-linux-amd64.tar.gz nginx-build
	GO111MODULE=on GOOS=darwin go build -ldflags '-s -w -X main.NginxBuildVersion=${VERSION}' -o nginx-build
	tar cvzf release/nginx-build-darwin-amd64.tar.gz nginx-build
	rm nginx-build

# ImageMagick and GD are required for ngx_small_light
build-example: nginx-build
	./nginx-build -c config/configure.example -m config/modules.cfg.example -d work -clear

check:
	GO111MODULE=on go test ./...

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
