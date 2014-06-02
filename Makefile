

nginx-build: *.go
	gom build -o nginx-build

build-cross:
	GOOS=linux GOARCH=amd64 gom build -o bin/linux/amd64/nginx-build
	GOOS=darwin GOARCH=amd64 gom build -o bin/darwin/amd64/nginx-build

build-example: nginx-build
	./nginx-build -c config/configure.example -m config/modules.cfg.example -d work -clear

bundle:
	gom install

check:
	gom test

clean:
	rm -rf nginx-build
	rm -rf bin/linux/amd64/nginx-build
	rm -rf bin/darwin/amd64/nginx-build
