package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/your-org/your-monorepo/pkg/userpb"
	"google.golang.org/grpc"
)

// UserService implementasyonu (boş)
type userServer struct {
	userpb.UnimplementedUserServiceServer
}

func (s *userServer) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.UserResponse, error) {
	// Burada kullanıcı oluşturma işlemi yapılacak
	return &userpb.UserResponse{}, nil
}

func (s *userServer) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.UserResponse, error) {
	// Burada kullanıcı getirme işlemi yapılacak
	return &userpb.UserResponse{}, nil
}

func (s *userServer) UpdateUser(ctx context.Context, req *userpb.UpdateUserRequest) (*userpb.UserResponse, error) {
	// Burada kullanıcı güncelleme işlemi yapılacak
	return &userpb.UserResponse{}, nil
}

func (s *userServer) DeleteUser(ctx context.Context, req *userpb.DeleteUserRequest) (*userpb.DeleteUserResponse, error) {
	// Burada kullanıcı silme işlemi yapılacak
	return &userpb.DeleteUserResponse{Success: true}, nil
}

func (s *userServer) ListUsers(ctx context.Context, req *userpb.ListUsersRequest) (*userpb.ListUsersResponse, error) {
	// Burada kullanıcı listeleme işlemi yapılacak
	return &userpb.ListUsersResponse{}, nil
}

func main() {
	// gRPC server başlat
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Port dinlenemedi: %v", err)
	}
	grpcServer := grpc.NewServer()
	userpb.RegisterUserServiceServer(grpcServer, &userServer{})
	fmt.Println("User gRPC servisi 50051 portunda başlatıldı...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("gRPC server başlatılamadı: %v", err)
	}
}
