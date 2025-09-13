build:
	npx tailwindcss -i ./css/input.css -o ./css/main.css --minify
	go-assets-builder schema.sql favicon.ico ini htm img css -o src/assets.go
	cd src; CGO_ENABLED=1 go build -o ../bin/sm -ldflags="-w -s" .

buildwin:
	npx tailwindcss -i ./css/input.css -o ./css/main.css --minify
	go-assets-builder schema.sql favicon.ico ini htm img css -o src/assets.go
	go-winres make
	cd src; CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 go build -o ../bin/sm.exe -ldflags='-w -s' .
	upx -9 bin/sm.exe

run:
	go-assets-builder schema.sql favicon.ico ini htm img css -o src/assets.go
	cd src; go run .

deps:
	go install github.com/jessevdk/go-assets-builder@latest
	go install github.com/tc-hib/go-winres@latest
	go get main/src
	npm install
	npx tailwindcss -i ./css/input.css -o ./css/main.css

clean:
	rm -r /bin
	rm css/main.css
	rm assets.go
	rm *.syso
	go clean
