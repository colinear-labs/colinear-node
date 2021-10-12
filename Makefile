SHELL = /bin/sh
.DEFAULT_GOAL := build

build:
	@mkdir -p node/bin
	cd node && (GOOS=linux GOARCH=arm go build -o bin/xnode-linux-arm; \
	GOOS=linux GOARCH=arm64 go build -o bin/xnode-linux-arm64; \
	GOOS=linux GOARCH=386 go build -o bin/xnode-linux-386; \
	GOOS=linux GOARCH=amd64 go build -o bin/xnode-linux-amd64; \
	)

clean:
	@rm -rf node/bin

