package main

import "fmt"

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
	//Use goriila mux -> sql -> Posgresdb
	fmt.Println("Hello World")
}
