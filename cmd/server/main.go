package main

import (
	"log"
	"net"

	v1 "github.com/idirall22/grpc_chat/api/pb"
	"github.com/idirall22/grpc_chat/pkg/service"

	"google.golang.org/grpc"
)

func main() {
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	srv := grpc.NewServer()
	v1.RegisterChatServiceServer(srv, service.NewService())
	log.Println("Server runing on localhost:8080")
	srv.Serve(l)
}
