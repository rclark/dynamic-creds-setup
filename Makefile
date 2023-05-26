clean:
	rm -rf dist/

build: clean
	GOOS=darwin GOARCH=amd64 go build -o dist/darwin/amd64/dynamic-creds-setup
	GOOS=darwin GOARCH=arm64 go build -o dist/darwin/arm64/dynamic-creds-setup
