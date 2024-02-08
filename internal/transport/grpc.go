package transport

import (
	"context"
	"log"
	"net"
	"strconv"

	"github.com/abdoroot/authentication-service/internal/auth"
	"github.com/abdoroot/authentication-service/internal/types"
	pb "github.com/abdoroot/authentication-service/proto"
	"google.golang.org/grpc"
)

type grpcTransport struct {
	pb.UnimplementedAuthenticationServiceServer
	srv        *auth.Auth
	server     *grpc.Server
	listenAddr string
}

func NewGRPCTransport(srv *auth.Auth, listenAddr string) *grpcTransport {
	return &grpcTransport{
		listenAddr: listenAddr,
		server:     grpc.NewServer(),
		srv:        srv,
	}
}

func (t *grpcTransport) Strart() error {
	log.Printf("grpc transport listen on port %v", t.listenAddr)
	l, err := net.Listen("tcp", t.listenAddr)
	if err != nil {
		log.Printf("Colud't listen on port %v", t.listenAddr)
		return err
	}
	pb.RegisterAuthenticationServiceServer(t.server, t)
	return t.server.Serve(l)
}

func (t *grpcTransport) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	resp := &pb.LoginResponse{}
	param := &types.LoginParam{
		Email:    req.Email,
		Password: req.Password,
	}
	user, err := t.srv.Login(ctx, param)
	if err != nil {
		return resp, err
	}
	userIdString := strconv.Itoa(user.ID)
	token, err := auth.GenerateToken(userIdString, user.Email)
	if err != nil {
		return resp, err
	}

	resp.AccessToken = token["access_token"]
	resp.RefreshToken = token["refresh_token"]
	return resp, nil
}

func (t *grpcTransport) SignUp(ctx context.Context, req *pb.SignUpRequest) (*pb.SignUpResponse, error) {
	resp := &pb.SignUpResponse{}
	param := &types.CreateUserParam{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}
	user, err := param.CreateUserFromParam()
	if err != nil {
		return resp, err
	}
	createdUser, err := t.srv.SignUp(context.Background(), user)
	_ = createdUser
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func (t *grpcTransport) Update(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	return &pb.UpdateResponse{}, nil
}

func (t *grpcTransport) UserProfile(ctx context.Context, req *pb.EmtpyRequest) (*pb.UserProfileResponse, error) {
	return &pb.UserProfileResponse{}, nil
}
