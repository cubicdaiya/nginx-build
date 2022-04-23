export GO111MODULE=on

nginx-build: *.go builder/*.go command/*.go configure/*.go module3rd/*.go openresty/*.go util/*.go
	go build -ldflags "-X main.NginxBuildVersion=`git rev-list HEAD -n1`" -o $@

build-example: nginx-build
	./nginx-build -c config/configure.example -m config/modules.cfg.example -d work -clear

check:
	go test ./...

fmt:
	go fmt ./...

clean:
	rm -rf nginx-build
