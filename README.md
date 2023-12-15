# Sidecar Application

This is a simple go web server that will host a POST endpoint (/validate) on port 8080.
The application is intended to run as a sidecar, but can be run stand alone.

# Docker commands
Docker files where created based off reference examples from 
The multistage docker file will create a compact image, since the go tool chain will not be included. 

- docker build -t zimmy71/sidecar:latest -f Dockerfile.multistage .
- docker run --publish 8071:8071 zimmy71/sidecar:latest
- docker push zimmy71/sidecar:latest