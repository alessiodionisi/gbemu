GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 CC=gcc go build -tags static -ldflags "-s -w" -o build/gbemu-darwin-amd64 -x cmd/gbemu/gbemu.go
