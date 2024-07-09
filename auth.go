package main

import (
	"context"
	"fmt"
	"log"
	pb "streakauth/grpc"
	"time"

	"github.com/golang-jwt/jwt"
)

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
	_, err := verifyToken(tokenString)
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

func verifyToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("error parsing claims")
	}

	username, ok := claims["username"].(string)
	if !ok {
		return "", fmt.Errorf("username claim not found or not a string")
	}

	return username, nil
}
