package main

import (
	"context"
	"log"
	"net"
	"sync"

	pb "grpc-hello/proto"	// go_package = "/proto" 설정 기준

	"google.golang.org/grpc"
	"github.com/google/uuid"
)

type helloServer struct {
	pb.UnimplementedHelloServiceServer
}

func (s *helloServer) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{Message: "Hello, " + req.Name + "!"}, nil
}

func (m *Matchmaker) Join(ctx context.Context, userId string) (*pb.JoinResponse, error) {
	m.mu.Lock()

	ch := make(chan *pb.JoinResponse)
	m.queue = append(m.queue, userId)
	m.waiters[userId] = ch

	if len(m.queue) >= 2 {
		// matching
		user1 := m.queue[0]
		user2 := m.queue[1]
		m.queue = m.queue[2:]

		matchId := uuid.New().String()
		res1 := &pb.JoinResponse{MatchId: matchId, OpponentId: user2, IsMatched: true}
		res2 := &pb.JoinResponse{MatchId: matchId, OpponentId: user1, IsMatched: true}

		ch1 := m.waiters[user1]
		ch2 := m.waiters[user2]
		go func() {
			ch1 <- res1
			ch2 <- res2
		}()
	}

	m.mu.Unlock()
	
	select {
	case res := <-ch:
		return res, nil
	case <-ctx.Done():
		m.mu.Lock()
		delete(m.waiters, userId)
		m.mu.Unlock()
		return nil, ctx.Err()
	}
}

type MatchServer struct {
	pb.UnimplementedMatchmakingServiceServer
	matchmaker *Matchmaker
}

type Matchmaker struct {
	queue []string
	waiters map[string]chan *pb.JoinResponse
	mu sync.Mutex
}

func NewMatchmaker() *Matchmaker {
	return &Matchmaker{
		queue: make([]string, 0),
		waiters: make(map[string]chan *pb.JoinResponse),
	}
}

func (s *MatchServer) JoinQueue(ctx context.Context, req *pb.JoinRequest) (*pb.JoinResponse, error) {
	return s.matchmaker.Join(ctx, req.UserId)
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	matchmaker := NewMatchmaker()
	pb.RegisterHelloServiceServer(grpcServer, &helloServer{})
	pb.RegisterMatchmakingServiceServer(grpcServer, &MatchServer{
		matchmaker: matchmaker,
	})

	log.Println("gRPC server is listening on port 50051...")
	grpcServer.Serve(lis)
}