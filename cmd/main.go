package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/brianvoe/gofakeit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	desc "github.com/ArturSaga/auth/api2/grpc/pkg/user_v1"
)

const grpcPort = 50051

type server struct {
	desc.UnimplementedUserApiServer
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterUserApiServer(s, &server{})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// CreateUser - публичный метод, который создает пользователя.
func (s *server) CreateUser(ctx context.Context, req *desc.CreateUserRequest) (*desc.CreateUserResponse, error) {
	log.Printf("Context: %+v", ctx)
	log.Printf("User id: %+v", req.GetInfo())

	return &desc.CreateUserResponse{}, nil
}

// GetUser - публичный метод, который позволяет получить данные пользователя.
func (s *server) GetUser(ctx context.Context, req *desc.GetUserRequest) (*desc.GetUserResponse, error) {
	log.Printf("User id: %+d", req.GetId())
	log.Printf("Context: %+v", ctx)
	password := gofakeit.Password(true, true, false, true, false, 9)

	return &desc.GetUserResponse{
		User: &desc.User{
			Id: req.GetId(),
			Info: &desc.UserInfo{
				Name:            gofakeit.Name(),
				Email:           gofakeit.Email(),
				Password:        password,
				PasswordConfirm: password,
				Role:            desc.Role_USER,
			},
			CreatedAt: timestamppb.New(gofakeit.Date()),
			UpdatedAt: timestamppb.New(gofakeit.Date()),
		},
	}, nil
}

// UpdateUser - публичный метод, который обновляет данные пользователя.
func (s *server) UpdateUser(ctx context.Context, req *desc.UpdateUserRequest) (*emptypb.Empty, error) {
	log.Printf("Context: %+v", ctx)
	log.Printf("User id: %+v", req.Info)
	return &emptypb.Empty{}, nil
}

// DeleteUser - публичный метод, который удаляет пользователя.
func (s *server) DeleteUser(ctx context.Context, req *desc.DeleteUserRequest) (*emptypb.Empty, error) {
	log.Printf("Context: %+v", ctx)
	log.Printf("User id: %+d", req.GetId())
	return &emptypb.Empty{}, nil
}
