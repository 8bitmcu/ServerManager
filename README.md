# Server Manager (SM)

SM is yet another web interface to manage a dedicated Assetto Corsa server. This project differs from others in which it uses presets to create events.


## Running on Linux

Using the included Dockerfile, the project can easily be run on linux using docker:
1. Clone the repository locally: `git clone git@github.com:8bitmcu/CorsaServer.git`
2. cd into the directory and build the dockerfile: `docker build --tag 'servermanager'`
3. Run the application in docker: `docker run 'servermanager'`
4. Access the web ui through the docker ip; you can get the IP by running `docker inspect ID` where ID is the container ID of the app found using `docker ps`

## Running on Windows

Coming soon


## Screenshots
