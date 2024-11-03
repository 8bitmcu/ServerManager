build:
	go-assets-builder schema.sql ini htm img css -o assets.go
	npx tailwindcss -i ./css/input.css -o ./css/main.css --minify
	CGO_ENABLED=1 go build -o bin/sm -ldflags="-w -s" main

buildwin:
	go-assets-builder schema.sql ini htm img css -o assets.go
	go-winres make
	npx tailwindcss -i ./css/input.css -o ./css/main.css --minify
	CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 go build -o bin/sm.exe -ldflags='-w -s' main
	upx bin/sm.exe

run:
	go-assets-builder schema.sql ini htm img css -o assets.go
	go run main

deps:
	go install github.com/jessevdk/go-assets-builder@latest
	go install github.com/tc-hib/go-winres@latest
	npm install
	npx tailwindcss -i ./css/input.css -o ./css/main.css

clean:
	rm css/main.css
	rm assets.go
	rm *.syso
	go clean
