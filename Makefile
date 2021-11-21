SHELL = /bin/sh
.DEFAULT_GOAL := release

abigen:
	cd node/processing/erc20eth/erc20abi && \
	solc --abi erc20.sol -o abi --overwrite && \
	abigen --abi=abi/ERC20.abi --pkg=erc20abi --type=ERC20 --out=erc20abi.go

build:
	@mkdir -p node/bin
	cd node && \
	GOOS=linux GOARCH=arm go build -o bin/x-node-linux-arm; \
	GOOS=linux GOARCH=arm64 go build -o bin/x-node-linux-arm64; \
	GOOS=linux GOARCH=386 go build -o bin/x-node-linux-386; \
	GOOS=linux GOARCH=amd64 go build -o bin/x-node-linux-amd64; \

clean:
	@rm -rf node/bin
	@rm -rf release

release: abigen build 
	@mkdir -p release
	@cp -r docker-compose.yml extnodes release
	@cp cli release/x-node
	@mkdir -p release/node-release
	@cp node-release.Dockerfile release/node-release/Dockerfile
	@cp node/bin/* release/node-release/
