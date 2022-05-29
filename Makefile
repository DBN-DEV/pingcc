all: gen-pb build

gen: install-gogrpc install-vtprotobuf gen-pb

build:
	go build -o pingcc main.go

gen-pb:
	mkdir ./pb && protoc --go-grpc_out=./pb --go_out=./pb --go-vtproto_out=./pb --go-vtproto_opt=features=marshal+unmarshal+size ./proto/*.proto

install-gogrpc:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

install-vtprotobuf:
	go install github.com/planetscale/vtprotobuf/cmd/protoc-gen-go-vtproto@latest
