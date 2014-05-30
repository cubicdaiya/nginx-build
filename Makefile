

nginx-build: *.go
	gom build

build-example: nginx-build
	./nginx-build -c config/configure.example -m config/modules.cfg.example -d work -clear

bundle:
	gom install

check:
	gom test

clean:
	rm -rf nginx-build
