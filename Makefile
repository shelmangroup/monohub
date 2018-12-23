all: check monohub

monohub: protoc
	go build .

protoc:
	protoc -I api/ --go_out=plugins=grpc:api monohub.proto

check: test

test:
	go test -v github.com/shelmangroup/monohub/...
