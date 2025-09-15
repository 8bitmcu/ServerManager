# Server Manager (SM)

SM is a very flexible web interface used to manage a dedicated Assetto Corsa server.

Features:

 - Supports both Windows and Linux
 - Single executable file, no external dependencies

## Running on Linux

### Prebuilt binaries

Use the pre-built binaries provided in the release tab.

### Dockerfile

The following instructions will build and run the application using the Dockerfile

1. Clone the repository locally
2. cd in the directory and build the Dockerfile
3. Run the application in Docker
4. Acceess the Web UI using http://localhost:3030

```sh
git clone https://github.com/8bitmcu/ServerManager.git
cd ServerManager
docker build . --tag 'servermanager'
docker run --network=host 'servermanager' -v /path/to/corsa:/corsa
```

### Docker compose

1. Clone the repository locally
2. cd into the directory and run `docker compose up`
3. Access the Web UI using https://localhost:443, http://localhost:80 or http://localhost:3030

## Running on Windows

Use the pre-built binaries provided in the release tab.

## Building the project

The project depends on tailwindcss to generate `/css/main.css`, and go-assets-builder to generate `assets.go`. npm is required for tailwindcss.
Run the following to install the dependencies

```sh
make deps
```

### Compiling for Linux

It is best to use the included Makefile to generate a build

```sh
make build
```

### Cross-compile for Windows

The makefile provides an easy to use method to create a build for Windows

```
make buildwin
```

## Screenshots
