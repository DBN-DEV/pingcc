


.PHONY: gen-pb
gen-pb:
	protoc --go-grpc_out=./pb --go_out=./pb --go-vtproto_out=./pb --go-vtproto_opt=features=marshal+unmarshal+size ./proto/*.proto
