clean:
	rm -rf dist/

build: clean
	GOOS=darwin GOARCH=amd64 go build -o dist/dynamic-creds-setup-darwin-amd64
	GOOS=darwin GOARCH=arm64 go build -o dist/dynamic-creds-setup-darwin-arm64
