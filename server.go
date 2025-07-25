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

type Matchmaker struct {
	queue []string
}

func (m *Matchmaker) Join(userId string) (isMatched bool, opponentId string) {
	m.queue = append(m.queue, userId)

	if len(m.queue) >= 2 {
		// matching
		user1 := m.queue[0]
		user2 := m.queue[1]
		m.queue = m.queue[2:]

		if userId == user1 {
			return true, user2
		}
		return true, user1
	}
	return false, ""
}

type matchServer struct {
	pb.UnimplementedMatchmakingServiceServer
	matchmaker *Matchmaker
}

func NewMatchmaker() *Matchmaker {
	return &Matchmaker{
		queue: make([]string, 0),
	}
}

func (s *matchServer) JoinQueue(ctx context.Context, req *pb.JoinRequest) (*pb.JoinResponse, error) {
	isMatched, opponentId := s.matchmaker.Join(req.UserId)
	if isMatched {
		return &pb.JoinResponse{
			MatchId: "match_1", 
			OpponentId: opponentId,
			IsMatched: true,
			}, nil		
	}
	return &pb.JoinResponse{
		MatchId: "", 
		OpponentId: "",
		IsMatched: false,
		}, nil	
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	matchmaker := NewMatchmaker()
	pb.RegisterHelloServiceServer(grpcServer, &helloServer{})
	pb.RegisterMatchmakingServiceServer(grpcServer, &matchServer{
		matchmaker: matchmaker,
	})

	log.Println("gRPC server is listening on port 50051...")
	grpcServer.Serve(lis)
}