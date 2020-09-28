package store

import (
	"context"

	v1 "github.com/idirall22/grpc_chat/api/pb"
)

// Store interface
type Store interface {
	Login(ctx context.Context, name string) (string, error)
	List(ctx context.Context) []*v1.User
}
