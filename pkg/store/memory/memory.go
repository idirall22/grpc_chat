package memory

import (
	"context"
	"sync"

	uuid "github.com/satori/go.uuid"

	v1 "github.com/idirall22/grpc_chat/api/pb"
)

//Store in memory store
type Store struct {
	users map[string]*v1.User
	mx    *sync.RWMutex
}

// NewMemoryStore create new memory store
func NewMemoryStore() *Store {
	return &Store{
		users: make(map[string]*v1.User),
		mx:    &sync.RWMutex{},
	}
}

// Login a user
func (s *Store) Login(ctx context.Context, name string) (string, error) {
	s.mx.Lock()
	defer s.mx.Unlock()
	id := uuid.NewV4().String()
	s.users[id] = &v1.User{
		Id:   id,
		Name: name,
	}
	return id, nil
}

// List all users
func (s Store) List(ctx context.Context) []*v1.User {
	s.mx.Lock()
	defer s.mx.Unlock()
	users := []*v1.User{}
	for _, user := range s.users {
		users = append(users, user)
	}
	return users
}
