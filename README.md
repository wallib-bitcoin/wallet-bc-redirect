# Redirect service

This service is a simple HTTP redirect that forwards incoming requests to a specified URL. It can handle GET, POST and PUT requests.

## Requirements

- Go 1.19 or later
- Docker
- Kubernetes

## Installation

To build and run this service, you will need to have Go and Docker installed.

### Clone the repository

git clone https://github.com/wallib-bitcoin/wallet-bc-redirect.git

### Build the service

cd redirect-service
go build

### Run the service

REDIRECT_URL=https://wallib.com ./redirect-service

## Usage

To test the service, you can use `curl` to send a GET or POST request to the service.

### GET request

    curl http://localhost:8080/path

### POST request

    curl -X POST -d '{"invoice":"123456","status":"pending","wallet_id":"123"}' http://localhost:8080/path


