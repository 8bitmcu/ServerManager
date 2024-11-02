build:
	go-assets-builder schema.sql ini htm img css -o assets.go
	export GIN_MODE=release
	CGO_ENABLED=1 go build -o bin/sm -ldflags="-w -s" main

buildwin:
	go-assets-builder schema.sql ini htm img css -o assets.go
	export GIN_MODE=release
	CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 go build -o bin/sm.exe -ldflags="-w -s" main
	upx bin/sm.exe

run:
	go-assets-builder schema.sql ini htm img css -o assets.go
	go run main

clean:
	go clean
