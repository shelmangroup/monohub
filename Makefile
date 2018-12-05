all: check monohub

monohub: 
	go build .

check: test

test:
	go test github.com/shelmangroup/monohub/...
