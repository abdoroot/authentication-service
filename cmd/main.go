package main

import (
	"log"
	"net"

	"github.com/abdoroot/authentication-service/internal/auth"
	pb "github.com/abdoroot/authentication-service/proto"
	_ "github.com/joho/godotenv/autoload"
	"google.golang.org/grpc"
)

//todo Handle token expiration and refresh if necessary
//todo refactor DB remove the glopal db

func main() {
	//connect to database
	db, err := auth.NewDB()
	if err != nil {
		log.Panic(err)
	}

	//migrate : use it when needed
	/*
		err = db.Migrate()
		if err != nil {
			log.Panic(err)
		}
	*/

	//create grpcserver && auth instance
	gs := grpc.NewServer(grpc.UnaryInterceptor(auth.AuthUInterceptor)) //grpc.UnaryInterceptor()//grpc middleware
	au := auth.NewAuth(db)                                             //auth handlers

	//regiter the grpc serve and the auth instance that implement ...
	pb.RegisterAuthenticationServiceServer(gs, au)

	//create net listener
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Printf("Colud't listen on port 8080")
	}

	//start the grpc server
	gs.Serve(l)
}
