package main

import (
	"context"
	"fmt"
	"log"
	pb "streakauth/grpc"
)

func (s *server) CheckAuthorized(ctx context.Context, in *pb.CheckAuthorizedReq) (*pb.CheckAuthorizedRes, error) {
	log.Printf("Received authorization check request: %v", in)

	username, err := verifyToken(in.AuthCode)
	if err != nil {
		fmt.Print("Invalid token")
		return &pb.CheckAuthorizedRes{Username: "", Authorized: false}, err
	}

	return &pb.CheckAuthorizedRes{Username: username, Authorized: true}, nil
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