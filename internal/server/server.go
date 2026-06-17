package server

import (
	"account/internal/mapper"
	"account/internal/model"
	"context"
	accountpb "github.com/matveevsa/contracts/account"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/rs/zerolog"
)

type Server struct {
	accountpb.UnimplementedAccountServer

	accountService AccountService
	lg             *zerolog.Logger
}

func New(accountService AccountService, lg *zerolog.Logger) *Server {
	return &Server{accountService: accountService, lg: lg}
}

type AccountService interface {
	CreateUser(ctx context.Context, user model.CreateUser) error
	GetUser(ctx context.Context, userID uint64) (model.User, error)
	GetUsers(ctx context.Context, offset, limit int) ([]model.User, error)
	DeleteUser(ctx context.Context, userID uint64) error
	UpdateUser(ctx context.Context, userID uint64, user model.UpdateUser) error
}

func (s *Server) CreateUser(ctx context.Context, req *accountpb.CreateUserRequest) (*emptypb.Empty, error) {
	if req.GetUser() == nil {
		return &emptypb.Empty{}, status.Error(codes.InvalidArgument, "user is required")
	}

	user := mapper.PbToUserCreate(req.GetUser())

	if err := s.accountService.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) GetUser(ctx context.Context, req *accountpb.GetUserRequest) (*accountpb.GetUserResponse, error) {
	user, err := s.accountService.GetUser(ctx, req.GetUserId())

	if err != nil {
		return nil, err
	}

	return &accountpb.GetUserResponse{User: mapper.UserToPb(user)}, nil
}

func (s *Server) GetUsers(ctx context.Context, req *accountpb.GetUsersRequest) (*accountpb.GetUsersResponse, error) {
	var offset = 0
	var limit = 10

	if req.GetPagination() != nil {
		offset = int(req.GetPagination().Offset)
		limit = int(req.GetPagination().GetLimit())
	}

	res, err := s.accountService.GetUsers(ctx, offset, limit)

	if err != nil {
		return nil, err
	}

	return &accountpb.GetUsersResponse{Users: mapper.UsersToPbs(res)}, nil
}

func (s *Server) DeleteUser(ctx context.Context, req *accountpb.DeleteUserRequest) (*emptypb.Empty, error) {
	if err := s.accountService.DeleteUser(ctx, req.GetUserId()); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) UpdateUser(ctx context.Context, req *accountpb.UpdateUserRequest) (*emptypb.Empty, error) {
	if req.GetUser() == nil {
		return nil, status.Error(codes.InvalidArgument, "user is required")
	}

	user := mapper.PbToUserUpdate(req.GetUser())

	if err := s.accountService.UpdateUser(ctx, req.GetUserId(), user); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
