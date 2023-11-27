
GOPATH:=$(shell go env GOPATH)
MODIFY=proto/

.PHONY: proto
proto:
    
	protoc --proto_path=${MODIFY} --micro_out=${MODIFY} --go_out=${MODIFY} order.proto
    

.PHONY: build
build: 

	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o order-service *.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker: build
	docker build . -t order-service:latest
