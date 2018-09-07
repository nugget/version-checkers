VER := $(shell git describe --tags)

.PHONY: preamble clean all

all: preamble build/macos/vc-zonedirector build/linux/vc-zonedirector build/freebsd/vc-zonedirector
	@echo Current builds:
	@find build -type f | xargs file

preamble:
	go get
	go generate

build/macos/vc-zonedirector:
	GOARCH=amd64 GOOS=darwin go build -o build/macos/vc-zonedirector

build/linux/vc-zonedirector:
	GOARCH=amd64 GOOS=linux go build -o build/linux/vc-zonedirector

build/freebsd/vc-zonedirector:
	GOARCH=amd64 GOOS=freebsd go build -o build/freebsd/vc-zonedirector

tarballs:
	cd build/macos && tar cjvf ../macos-vc-zonedirector-$(VER).tar.bz2 vc-zonedirector
	cd build/linux && tar cjvf ../linux-vc-zonedirector-$(VER).tar.bz2 vc-zonedirector
	cd build/freebsd && tar cjvf ../freebsd-vc-zonedirector-$(VER).tar.bz2 vc-zonedirector

install:
	cp -rp build/macos/vc-zonedirector $(HOME)/bin/

zabbix: build/freebsd/vc-zonedirector
	chmod 755 build/freebsd/vc-zonedirector
	scp -rp build/freebsd/vc-zonedirector zabbix.nuggethaus.net:/usr/local/etc/zabbix34/zabbix/externalscripts

clean:
	rm -rf build
