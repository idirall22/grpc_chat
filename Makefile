gen:
	protoc \
	--proto_path=api/proto \
	--proto_path=${GOPATH}/src/github.com/gogo/protobuf/protobuf \
	--proto_path=${GOPATH}/src \
	--gofast_out=plugins=grpc:api/pb/ api/proto/*.proto \

server: 
	go run cmd/server/main.go

client1: 
	go run cmd/client/main.go -user 2

client2: 
	go run cmd/client/main.go -user 1