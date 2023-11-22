package auth

import (
	"context"
	"log"
	"testing"

	pb "github.com/abdoroot/authentication-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

var (
	Atoken string = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDAyMTU4NDgsInVzZXJfZW1haWwiOiJhYmQuMjAwOTMwQGdtYWlsLmNvbSIsInVzZXJfaWQiOiIyMCJ9.XyS44pR-tom02j5ByjXzloGOjKfKF8qaDSrTHpoQX6s"

	rToken string = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDA3Njk3NDIsInVzZXJfaWQiOiIyMCJ9.aG5VNBnhJwhScs-S-k0ztHe9T2o7_gK4OfFw4n7j68w"

	Server string = "localhost:8080"

	Conn *grpc.ClientConn
)

func init() {
	var err error
	Conn, err = grpc.Dial(Server, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("fail to connect to %v with error %v", Server, err)
	}
}

func TestUpdate(t *testing.T) {
	c := pb.NewAuthenticationServiceClient(Conn)
	// Create a context with the token value using context.WithValue
	ctx := metadata.NewOutgoingContext(
		context.Background(),
		metadata.Pairs("access_token", Atoken),
	)
	resp, err := c.Update(ctx, &pb.UpdateRequest{
		Name: "Abdelhadi Mohammed",
	})
	if err != nil {
		t.Error(err)
	}
	_ = resp
}

func TestGetProfile(t *testing.T) {
	c := pb.NewAuthenticationServiceClient(Conn)
	// Create a context with the token value using context.WithValue
	ctx := metadata.NewOutgoingContext(
		context.Background(),
		metadata.Pairs("access_token", Atoken),
	)
	resp, err := c.UserProfile(ctx, &pb.EmtpyRequest{})
	if err != nil {
		t.Error(err)
	}
	if resp.Name == "" && resp.Email == "" {
		t.Error("error getting data from db")
	}
}

func TestAcessToken(t *testing.T) {
	c := pb.NewAuthenticationServiceClient(Conn)
	// Create a context with the token value using context.WithValue
	ctx := metadata.NewOutgoingContext(
		context.Background(),
		metadata.Pairs("access_token", Atoken, "refresh_token", rToken),
	)
	resp, err := c.UserProfile(ctx, &pb.EmtpyRequest{})
	if err != nil {
		t.Error(err)
	}
	_ = resp
}
