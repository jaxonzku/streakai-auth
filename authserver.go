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
var currentUsers = []string{}

func (s *server) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	// Simulate some processing time

	log.Printf("Received Login request: %v", in)

	userExist := isUserRegistered(in.Username)
	var tokenString string
	if userExist {

		if registeredUsers[in.Username] != in.Password {

			fmt.Errorf("No username found")
			return &pb.LoginResponse{Token: ""}, nil

		} else {
			tokenString, err := CreateToken(in.Username)

			if err != nil {
				fmt.Errorf("Error generating token", err)
			}

			return &pb.LoginResponse{Token: tokenString}, nil

		}

	}
	return &pb.LoginResponse{Token: tokenString}, nil

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
