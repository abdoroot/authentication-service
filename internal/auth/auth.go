package auth

import (
	"context"
	"log"
	"net/mail"

	pb "github.com/abdoroot/authentication-service/proto"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type auth struct {
	dbi *sqlx.DB //datbase package instant
	pb.UnimplementedAuthenticationServiceServer
}

func NewAuth(dbi *sqlx.DB) *auth {
	return &auth{
		dbi: dbi,
	}
}

func (a auth) SignUp(ctx context.Context, req *pb.SignUpRequest) (*pb.SignUpResponse, error) {
	if !a.validateSignUp(req) {
		return nil, status.Errorf(codes.InvalidArgument, "Input validation error")
	}
	err := DbInsert(req)
	if err != nil {
		return nil, err
	}
	return &pb.SignUpResponse{}, nil
}

func (a auth) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	tmp, err := DbLogin(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "error username or password")
	}

	t, rt := tmp["access_token"], tmp["refresh_token"]
	return &pb.LoginResponse{
		AccessToken:  t,
		RefreshToken: rt,
	}, nil
}

func (a auth) Update(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	token := parseToken(ctx)
	claims, auth := IsUserAuthorizedWithClaim(token)
	if auth {
		//update opration
		//todo validate inputs
		log.Println("Auth user")
		if err := DbUpdate(req, claims); err != nil {
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
		resp, err := DbGetProfile(claims)
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
	if t, ok := md["access_token"]; ok {
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

func AuthUInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Println("Unary interceptor/middleware invoked,", info)
	return handler(ctx, req)
}
