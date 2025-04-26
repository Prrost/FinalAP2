package handlers

import (
	"context"
	userpb "github.com/Prrost/protoFinalAP2/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
	"user-service/config"
	"user-service/domain"
	"user-service/useCase"
)

type UserServer struct {
	userpb.UnimplementedUserServiceServer
	cfg *config.Config
	uc  *useCase.UseCase
}

func NewUserServer(cfg *config.Config, uc *useCase.UseCase) *UserServer {
	return &UserServer{
		cfg: cfg,
		uc:  uc,
	}
}

func (s *UserServer) RegisterUser(ctx context.Context, req *userpb.RegisterRequest) (*userpb.RegisterResponse, error) {
	var user, userOut domain.User
	var err error

	user.Email = req.GetEmail()
	user.Password = req.GetPassword()
	user.IsAdmin = req.GetIsAdmin()

	if user.IsAdmin {
		if user.Password == "" {
			return nil, status.Error(codes.InvalidArgument, "password must be provided for admin")
		}
		userOut, err = s.uc.CreateUserAdmin(user)
	} else {
		userOut, err = s.uc.CreateUser(user)
		if err != nil {
			return nil, err
		}
	}

	return &userpb.RegisterResponse{
		Id:      int64(userOut.ID),
		Message: "User created successfully",
	}, nil
}

func (s *UserServer) AuthenticateUser(ctx context.Context, req *userpb.AuthRequest) (*userpb.AuthResponse, error) {
	var user domain.User

	user.Email = req.GetEmail()
	user.Password = req.GetPassword()

	token, err := s.uc.LoginUser(user)
	if err != nil {
		return nil, err
	}

	return &userpb.AuthResponse{
		Token:   token,
		Message: "Authentication successful",
	}, nil
}

func (s *UserServer) GetUserInfo(ctx context.Context, req *userpb.UserInfoRequest) (*userpb.UserInfoResponse, error) {
	id := req.GetId()
	email := req.GetEmail()
	var user domain.User
	var err error
	var idString string

	if id == "" && email == "" {
		return nil, status.Error(codes.InvalidArgument, "id or email must be provided")
	}

	if id != "" {
		user, err = s.uc.GetUserByID(id)
		if err != nil {
			return nil, err
		}

	}

	if email != "" {
		user, err = s.uc.GetUserByEmail(email)
		if err != nil {
			return nil, err
		}
	}

	idString = strconv.Itoa(int(user.ID))

	return &userpb.UserInfoResponse{
		Id:      idString,
		Email:   user.Email,
		IsAdmin: user.IsAdmin,
	}, nil
}
