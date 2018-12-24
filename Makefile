all: check monohub

monohub: protoc
	go build .

protoc:
	protoc \
		-I api/ \
		-I vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/ \
		--go_out=plugins=grpc:api \
		monohub.proto
	protoc \
		-I api/ \
		-I vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/ \
		--grpc-gateway_out=logtostderr=true:api \
		monohub.proto

check: test

test:
	go test -v github.com/shelmangroup/monohub/...
