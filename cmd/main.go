package main

import (
	"log"
	"net"

	"github.com/abdoroot/authentication-service/internal/auth"
	"github.com/abdoroot/authentication-service/internal/database"
	pb "github.com/abdoroot/authentication-service/proto"
	_ "github.com/joho/godotenv/autoload"
	"google.golang.org/grpc"
)

/*
   User Registration:
       Allow users to register by providing necessary information like name, email, and password.

   User Login:
       Implement a login functionality where users can authenticate with their email and password.

   Token Generation:
       Use JWT (JSON Web Tokens) to generate tokens upon successful authentication.
       Include the user's information in the token payload.

   Token Validation:
       Create a mechanism to validate incoming tokens to ensure their authenticity.
       Handle token expiration and refresh if necessary.

   Password Hashing:
       Ensure the security of user passwords by using a strong hashing algorithm (e.g., bcrypt).

   User Profile Management:
       Allow users to update their profiles, including changing passwords.

   Role-Based Access Control (Optional):
       Implement role-based access control to restrict access to certain functionalities based on user roles.
*/

func main() {
	//connect to database
	db, err := database.NewDB()
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
	gs := grpc.NewServer()
	au := auth.NewAuth(db) //auth handlers

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
