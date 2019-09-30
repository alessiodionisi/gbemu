GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc go build -tags static -ldflags "-s -w" -o build/gbemu-windows-amd64.exe -x cmd/gbemu/gbemu.go
