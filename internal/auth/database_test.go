package auth

import (
	"context"
	"testing"

	pb "github.com/abdoroot/authentication-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

var Token string = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDAxNTkxNzEsInVzZXJfZW1haWwiOiJhYmQuMjAwOTMwQGdtYWlsLmNvbSIsInVzZXJfaWQiOiIxMCJ9.e4RFUNA0HyPPy8MmGiJOdjgM9By1ySe-Mpy8uxr7luY"

func TestUpdate(t *testing.T) {
	server := "localhost:8080"
	conn, err := grpc.Dial(server, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Errorf("fail to connect to %v with error %v", server, err)
		return
	}
	c := pb.NewAuthenticationServiceClient(conn)

	// Create a context with the token value using context.WithValue
	ctx := metadata.NewOutgoingContext(
		context.Background(),
		metadata.Pairs("token", Token),
	)

	resp, err := c.Update(ctx, &pb.UpdateRequest{})
	if err != nil {
		t.Error(err)
	}
	t.Error(resp)
}
