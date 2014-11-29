

nginx-build: *.go
	gom build -o nginx-build

build-cross:
	GOOS=linux GOARCH=amd64 gom build -o bin/linux/amd64/nginx-build
	GOOS=darwin GOARCH=amd64 gom build -o bin/darwin/amd64/nginx-build

dist: build-cross
	cd bin/linux/amd64/ && tar zcvf nginx-build-linux-amd64.tar.gz nginx-build
	cd bin/darwin/amd64/ && tar zcvf nginx-build-darwin-amd64.tar.gz nginx-build

# ImageMagick and GD are required
build-example: nginx-build
	./nginx-build -c config/configure.example -m config/modules.cfg.example -d work -clear

bundle:
	gom install

check:
	gom test

fmt:
	go fmt ./...

.PHONY: doc
doc:
	cd doc; make man

clean:
	rm -rf nginx-build
	rm -rf bin/linux/amd64/nginx-build*
	rm -rf bin/darwin/amd64/nginx-build*
