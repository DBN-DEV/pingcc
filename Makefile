all: build

build:
	go build -o bin/pingcc main.go

fmt:
	goimports -local github.com/DBN-DEV/pingcc -l -w .

gen-mocks:
	mockery --all
