all: check monohub

monohub: 
	go build .

check: test

test:
	go test -v github.com/shelmangroup/monohub/...
