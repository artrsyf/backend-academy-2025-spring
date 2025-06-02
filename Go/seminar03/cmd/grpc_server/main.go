package main

import (
	"context"
	"log"
	"net"
	"time"

	pb "example.com/seminar03/internal/api/proto/service/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	pb.UnimplementedEchoServiceServer
}

func (s *server) Echo(ctx context.Context, req *pb.EchoRequest) (*pb.EchoResponse, error) {
	return &pb.EchoResponse{
		Message: "Echo: " + req.Message,
		Time:    time.Now().Format(time.RFC1123Z),
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	pb.RegisterEchoServiceServer(grpcServer, &server{})
	log.Println("Starting server on :50051")
	grpcServer.Serve(lis)
}
