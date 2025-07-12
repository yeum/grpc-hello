package main

import (
	"context"
	"log"
	"net"

	pb "grpc-hello/proto"	// go_package = "/proto" 설정 기준

	"google.golang.org/grpc"
)

type helloServer struct {
	pb.UnimplementedHelloServiceServer
}

func (s *helloServer) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{Message: "Hello, " + req.Name + "!"}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterHelloServiceServer(grpcServer, &helloServer{})

	log.Println("gRPC server is listening on port 50051...")
	grpcServer.Serve(lis)
}