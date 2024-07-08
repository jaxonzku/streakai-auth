package main

import (
	"context"
	"fmt"
	"log"
	"net"
	pb "streakauth/grpc"
	"time"

	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedStreakAiServiceServer
}

var secretKey = []byte("secret-key")
var registeredUsers = map[string]string{}
var loggedinUsers = []string{}

var currentUsers = []string{}

func (s *server) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	// Simulate some processing time

	log.Printf("Received Login request: %v", in)

	userExist := isUserRegistered(in.Username)
	if userExist {

		if registeredUsers[in.Username] != in.Password {

			fmt.Errorf("No username found")
			return &pb.LoginResponse{Token: ""}, nil

		} else {
			tokenString, err := CreateToken(in.Username)
			loggedinUsers = append(loggedinUsers, in.Username)
			fmt.Println("loggedinUsers", loggedinUsers)

			if err != nil {
				fmt.Errorf("Error generating token", err)
			}

			return &pb.LoginResponse{Token: tokenString}, nil

		}

	}
	return nil, nil

}

func (s *server) LogOut(ctx context.Context, in *pb.LogOutRequest) (*pb.LogOutResponse, error) {
	// Simulate some processing time

	log.Printf("Received Logout request: %v", in)

	tokenString := in.AuthCode
	if tokenString == "" {
		fmt.Print("Missing authorization code")
		return nil, nil
	}
	err := verifyToken(tokenString)
	if err != nil {
		fmt.Print("Invalid token")
		return nil, nil
	}

	removeFromLoggedIn(in.Username)
	loggedinUsers = append(loggedinUsers, in.Username)

	return &pb.LogOutResponse{Status: "Logged Out"}, nil

}

func (s *server) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	// Simulate some processing time

	log.Printf("Received Register request: %v", in)

	if !isUserRegistered(in.Username) {
		registeredUsers[in.Username] = in.Password
	}

	fmt.Println("registeredUsers", registeredUsers)

	return &pb.RegisterResponse{Status: "success"}, nil
}

func main() {

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen on port 50051: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterStreakAiServiceServer(s, &server{})
	log.Printf("gRPC server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}

func isUserRegistered(username string) bool {
	_, exists := registeredUsers[username]
	return exists
}

func removeFromLoggedIn(username string) {
	for i, u := range loggedinUsers {
		if u == username {
			loggedinUsers = append(loggedinUsers[:i], loggedinUsers[i+1:]...)
		}
	}
}

func CreateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func verifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}
