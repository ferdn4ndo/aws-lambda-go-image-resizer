#!/bin/bash

set -o pipefail

go mod download github.com/aws/aws-lambda-go
go mod download github.com/awslabs/aws-lambda-go-api-proxy
go mod download github.com/disintegration/imaging
go mod download github.com/gorilla/mux
go mod download github.com/BurntSushi/toml

go mod vendor

GOOS=linux GOARCH=amd64 go build -o build/main main.go

zip build/main.zip build/main
