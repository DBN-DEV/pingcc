all: gen-pb build

build:
	go build -o pingcc main.go

gen-pb:
	protoc --go-grpc_out=./pb --go_out=./pb --go-vtproto_out=./pb --go-vtproto_opt=features=marshal+unmarshal+size ./proto/*.proto

install-vtprotobuf:
	go install github.com/planetscale/vtprotobuf/cmd/protoc-gen-go-vtproto@latest
