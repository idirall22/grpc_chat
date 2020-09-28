package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	v1 "github.com/idirall22/grpc_chat/api/pb"

	"google.golang.org/grpc"
)

var id string
var toUserID string

func main() {
	cc, err := grpc.Dial(":8080", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer cc.Close()

	client := v1.NewChatServiceClient(cc)
	stream, err := client.Chat(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	err = stream.Send(&v1.ChatStreamRequest{
		Type: &v1.ChatStreamRequest_LoginRequest{
			LoginRequest: &v1.LoginRequest{Name: "idir"},
		},
	})

	if err != nil {
		log.Fatal(err)
	}
	m, err := stream.Recv()
	if err != nil {
		log.Fatal(err)
	}
	res := (m.Type).(*v1.ChatStreamResponse_LoginResponse)
	id = res.LoginResponse.Id

	go func() {
		for {
			m, err := stream.Recv()
			if err != nil {
				log.Fatal(err)
			}
			res := (m.Type).(*v1.ChatStreamResponse_MessageResponse)
			fmt.Println("-------------------------")
			fmt.Println(res.MessageResponse.Message)
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)

	// cmd:
	// user:id example: user:1
	// users
	// msg:string ex: msg:hello
	for scanner.Scan() {
		in := scanner.Text()
		cmd := strings.Split(in, ":")
		if len(cmd) <= 1 {
			if cmd[0] == "users" {
				resList, err := client.List(context.Background(), &v1.ListUsersRequest{})
				if err != nil {
					log.Fatal(err)
				}
				for _, user := range resList.Users {
					fmt.Println(user)
				}
				continue
			}
			fmt.Println("Please try again")
			continue
		}
		switch strings.ToLower(cmd[0]) {
		case "user":
			toUserID = cmd[1]
			break
		case "msg":
			err := stream.Send(&v1.ChatStreamRequest{
				Type: &v1.ChatStreamRequest_MessageRequest{
					MessageRequest: &v1.MessageRequest{
						Message: &v1.MessageChat{
							FromUserId: id,
							ToUserId:   toUserID,
							Message:    in,
						},
					},
				},
			})
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
