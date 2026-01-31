VERSION ?= dev-build

.PHONY: all build buildwin run deps package clean

all: deps clean buildwin build


build:
	mkdir -p bin/
	npx @tailwindcss/cli -i ./css/input.css -o ./css/main.css --minify
	go-assets-builder schema.sql favicon.ico ini htm img css -o src/assets.go
	cd src; CGO_ENABLED=1 go build -o ../bin/sm_linux -ldflags="-w -s -X 'main.Version=$(VERSION)'" .

buildwin:
	mkdir -p bin/
	npx @tailwindcss/cli -i ./css/input.css -o ./css/main.css --minify
	go-assets-builder schema.sql favicon.ico ini htm img css -o src/assets.go
	go-winres make
	cd src; CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 go build -o ../bin/sm_win.exe -ldflags="-w -s -X 'main.Version=$(VERSION)'" .

run:
	go-assets-builder schema.sql favicon.ico ini htm img css -o src/assets.go
	cd src; go run .

deps:
	go install github.com/jessevdk/go-assets-builder@latest
	go install github.com/tc-hib/go-winres@latest
	go get main/src
	npm install
	npx tailwindcss -i ./css/input.css -o ./css/main.css

package: clean buildwin build

clean:
	rm -r bin/
	rm -f css/main.css
	rm -f src/assets.go
	rm -f *.syso
	go clean
