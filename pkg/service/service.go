package service

import (
	"context"
	"fmt"
	"log"

	"github.com/idirall22/grpc_chat/pkg/store/memory"

	v1 "github.com/idirall22/grpc_chat/api/pb"
	"github.com/idirall22/grpc_chat/pkg/store"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserConnection struct
type UserConnection struct {
	stream   v1.ChatService_ChatServer
	messages chan *v1.MessageChat
}

// Service create new chat service
type Service struct {
	store store.Store
	// messages         []chan *v1.MessageChat
	UsersConnections map[string]*UserConnection
}

// NewService new service
func NewService() *Service {
	return &Service{
		store:            memory.NewMemoryStore(),
		UsersConnections: make(map[string]*UserConnection),
	}
}

// Chat endpoint
func (s *Service) Chat(stream v1.ChatService_ChatServer) error {
	req, err := stream.Recv()
	if err != nil {
		return status.Errorf(codes.Internal, err.Error())
	}

	loginReq := (req.Type).(*v1.ChatStreamRequest_LoginRequest)
	id, _ := s.store.Login(stream.Context(), loginReq.LoginRequest.Name)
	s.UsersConnections[id] = &UserConnection{stream, make(chan *v1.MessageChat, 64)}

	stream.Send(&v1.ChatStreamResponse{
		Type: &v1.ChatStreamResponse_LoginResponse{
			LoginResponse: &v1.LoginResponse{Id: id},
		},
	})

	msgChan := s.UsersConnections[id].messages

	go func() {
		for {
			select {
			case m := <-msgChan:
				err := stream.Send(&v1.ChatStreamResponse{
					Type: &v1.ChatStreamResponse_MessageResponse{
						MessageResponse: &v1.MessageResponse{
							Message: m,
						},
					},
				})
				if err != nil {
					log.Fatal(err)
				}

			}
		}
	}()

	for {
		req, err := stream.Recv()
		if err != nil {
			log.Fatal(err)
		}
		messageReq := (req.Type).(*v1.ChatStreamRequest_MessageRequest)
		s.UsersConnections[messageReq.MessageRequest.Message.ToUserId].messages <- messageReq.MessageRequest.Message
		fmt.Println("New Message Received", messageReq.MessageRequest.Message)
	}

	return nil
}

// List endpoint
func (s *Service) List(ctx context.Context, req *v1.ListUsersRequest) (*v1.ListUsersResponse, error) {
	users := s.store.List(ctx)
	return &v1.ListUsersResponse{
		Users: users,
	}, nil
}
