build-darwin-amd64:
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 CC=gcc go build -tags static -ldflags "-s -w" -o bin/darwin-amd64/gbemu -x cmd/gbemu/main.go

build-linux-amd64:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=1 CC=gcc go build -tags static -ldflags "-s -w" -o bin/linux-amd64/gbemu -x cmd/gbemu/main.go

build-windows-amd64:
	GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc go build -tags static -ldflags "-s -w" -o bin/windows-amd64/gbemu.exe -x cmd/gbemu/main.go

build-all: build-darwin-amd64 build-linux-amd64 build-windows-amd64
