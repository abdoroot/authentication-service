package auth

import (
	"context"
	"net/mail"

	pb "github.com/abdoroot/authentication-service/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type auth struct {
	dbi *DB //datbase package instant
	pb.UnimplementedAuthenticationServiceServer
}

func NewAuth(dbi *DB) *auth {
	return &auth{
		dbi: dbi,
	}
}

func (a auth) SignUp(ctx context.Context, req *pb.SignUpRequest) (*pb.SignUpResponse, error) {
	if !a.validate(req) {
		return nil, status.Errorf(codes.InvalidArgument, "Input validation error")
	}
	err := a.dbi.Insert(req)
	if err != nil {
		return nil, err
	}
	return &pb.SignUpResponse{}, nil
}

func (a auth) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	token, err := a.dbi.Login(req)
	if err != nil {
		return nil, err
	}

	return &pb.LoginResponse{
		Token: token,
	}, nil
}

func (a auth) Update(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	auth := IsUserAuthorized(ctx)
	_ = auth
	return &pb.UpdateResponse{}, nil
}

func (a auth) UserProfile(ctx context.Context, req *pb.EmtpyRequest) (*pb.UserProfileResponse, error) {
	auth := IsUserAuthorized(ctx)
	_ = auth
	return &pb.UserProfileResponse{}, nil
}

func (a auth) validate(req *pb.SignUpRequest) bool {
	if req.Name == "" || req.Email == "" || req.Password == "" {
		return false
	}
	if _, err := mail.ParseAddress(req.Email); err != nil {
		return false
	}
	return true
}
