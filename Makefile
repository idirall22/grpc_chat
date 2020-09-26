gen:
	protoc \
	--proto_path=api/proto \
	--proto_path=${GOPATH}/src/github.com/gogo/protobuf/protobuf \
	--proto_path=${GOPATH}/src \
	--gofast_out=plugins=grpc:api/pb/ api/proto/*.proto \

server: 
	go run cmd/main.go