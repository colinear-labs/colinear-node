SHELL = /bin/sh
.DEFAULT_GOAL := release

abigen:
	cd node/intents/erc20eth/erc20abi && \
	solc --abi erc20.sol -o abi --overwrite && \
	abigen --abi=abi/ERC20.abi --pkg=erc20abi --type=ERC20 --out=erc20abi.go

build:
	@mkdir -p node/bin
	cd node && \
	GOOS=linux GOARCH=arm go build -o bin/xnode-linux-arm; \
	GOOS=linux GOARCH=arm64 go build -o bin/xnode-linux-arm64; \
	GOOS=linux GOARCH=386 go build -o bin/xnode-linux-386; \
	GOOS=linux GOARCH=amd64 go build -o bin/xnode-linux-amd64; \

clean:
	@rm -rf node/bin
	@rm -rf release

release: abigen build 
	@mkdir -p release
	@cp -r cli docker-compose.yml extnodes release
	@mkdir -p release/node
	@cp node-release.Dockerfile release/node/Dockerfile
