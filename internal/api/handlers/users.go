package handlers

import (
	"context"
	"fmt"
	"os"

	"github.com/grpc_sqlc/database"
	"github.com/grpc_sqlc/db"
	mainapi "github.com/grpc_sqlc/proto/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	mainapi.UnimplementedUserServiceServer
}

func (s *Server) CreateUser(ctx context.Context, req *mainapi.UserRequest) (*mainapi.UserResponse, error) {
	// Implementation for creating a user
	conn, err := database.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, status.Error(codes.Aborted, fmt.Sprintf("Unable to connect to database: %v", err))
	}
	defer conn.Close(ctx)
	q := db.New(conn)
	_, err = q.CreateUser(ctx, db.CreateUserParams{
		FirstName:   req.GetFirstName(),
		LastName:    req.GetLastName(),
		PhoneNumber: req.GetPhoneNumber(),
	})

	if err != nil {
		return nil, status.Error(codes.Aborted, fmt.Sprintf("Unable to create user: %v", err))
	}

	return &mainapi.UserResponse{Message: "user created", Status: true}, nil
}

func (s *Server) GetUser(ctx context.Context, req *mainapi.UserIdRequest) (*mainapi.UserResponse, error) {
	// Implementation for getting a user
	conn, err := database.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, status.Error(codes.Aborted, fmt.Sprintf("Unable to connect to database: %v", err))
	}
	defer conn.Close(ctx)

	q := db.New(conn)
	_, err = q.GetUser(ctx, req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("User not found: %v", err))
	}
	return &mainapi.UserResponse{Message: "user found", Status: true}, nil
}

func (s *Server) UpdateUser(ctx context.Context, req *mainapi.UserUpdateRequest) (*mainapi.UserResponse, error) {
	// Implementation for updating a user
	conn, err := database.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, status.Error(codes.Aborted, fmt.Sprintf("Unable to connect to database: %v", err))
	}
	defer conn.Close(ctx)

	q := db.New(conn)
	err = q.UpdateUser(ctx, db.UpdateUserParams{
		ID:          req.GetUserId(),
		FirstName:   req.GetFirstName(),
		LastName:    req.GetLastName(),
		PhoneNumber: req.GetPhoneNumber(),
	})
	if err != nil {
		return nil, status.Error(codes.Aborted, fmt.Sprintf("Unable to update user: %v", err))
	}
	return &mainapi.UserResponse{Message: "user updated", Status: true}, nil
}

func (s *Server) DeactivateUser(ctx context.Context, req *mainapi.UserIdRequest) (*mainapi.UserResponse, error) {
	// Implementation for deactivating a user
	conn, err := database.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, status.Error(codes.Aborted, fmt.Sprintf("Unable to connect to database: %v", err))
	}
	defer conn.Close(ctx)

	q := db.New(conn)
	err = q.DeactivateUser(ctx, req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.Aborted, fmt.Sprintf("Unable to deactivate user: %v", err))
	}
	return &mainapi.UserResponse{Message: "user deactivated", Status: true}, nil
}
