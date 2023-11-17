package auth

import (
	"context"
	"log"
	"net/mail"

	pb "github.com/abdoroot/authentication-service/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
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
	if !a.validateSignUp(req) {
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
		return nil, status.Error(codes.InvalidArgument, "error username or password")
	}

	return &pb.LoginResponse{
		Token: token,
	}, nil
}

func (a auth) Update(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	token := parseToken(ctx)
	claims, auth := IsUserAuthorizedWithClaim(token)
	if auth {
		//update opration
		//todo validate inputs
		log.Println("Auth user")
		if err := a.dbi.Update(req, claims); err != nil {
			return &pb.UpdateResponse{}, status.Error(codes.Internal, "update error!")
		}
		//updated
		return &pb.UpdateResponse{}, status.Error(codes.OK, "updated Succesfully")
	}
	return &pb.UpdateResponse{}, status.Error(codes.Unauthenticated, "Unauthenticated")
}

func (a auth) UserProfile(ctx context.Context, req *pb.EmtpyRequest) (*pb.UserProfileResponse, error) {
	token := parseToken(ctx)
	claims, auth := IsUserAuthorizedWithClaim(token)
	if auth {
		//update opration
		//todo validate inputs
		resp, err := a.dbi.GetProfile(claims)
		if err != nil {
			return &pb.UserProfileResponse{}, status.Error(codes.Internal, "get profile error!")
		}
		//retrive
		return &pb.UserProfileResponse{
			Name:  resp.Name,
			Email: resp.Email,
		}, status.Error(codes.OK, "Retrieve Succesfully")
	}
	return &pb.UserProfileResponse{}, status.Error(codes.Unauthenticated, "Unauthenticated")
}

func parseToken(ctx context.Context) string {
	md, _ := metadata.FromIncomingContext(ctx)
	token := ""
	if t, ok := md["token"]; ok {
		token = t[0] //slice
	}
	return token
}

func (a auth) validateSignUp(req *pb.SignUpRequest) bool {
	if req.Name == "" || req.Email == "" || req.Password == "" {
		return false
	}
	if _, err := mail.ParseAddress(req.Email); err != nil {
		return false
	}
	return true
}
