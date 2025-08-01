package main

import (
	"context"
	"log"
	"net"
	"sync"
	"math/rand"

	pb "grpc-hello/proto"	// go_package = "/proto" 설정 기준

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
    "google.golang.org/grpc/codes"
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

		go func() {
			battleResult := StartBattle(matchId, user1, user2)
			m.mu.Lock()
			m.battleResults[matchId] = battleResult
			m.mu.Unlock()
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
	battleResults map[string]*BattleResult
	mu sync.Mutex
}

func NewMatchmaker() *Matchmaker {
	return &Matchmaker{
		queue: make([]string, 0),
		waiters: make(map[string]chan *pb.JoinResponse),
		battleResults: make(map[string]*BattleResult),
	}
}

func (s *MatchServer) JoinQueue(ctx context.Context, req *pb.JoinRequest) (*pb.JoinResponse, error) {
	return s.matchmaker.Join(ctx, req.UserId)
}

type Player struct {
	UserId string
	Hp int
}

type BattleResult struct {
	MatchId string
	Player1 string
	Player2 string
	Winner string
	Turns int
}

func StartBattle(matchId, user1, user2 string) *BattleResult {
	turns := 0
	players := []Player{
		{UserId: user1, Hp: 100}, 
		{UserId: user2, Hp: 100},
	}
	isNeededShuffle := rand.Intn(2)
	if isNeededShuffle == 1 {
		players[0], players[1] = players[1], players[0]
	}
	var attacker Player
	for players[0].Hp > 0 && players[1].Hp > 0 {
		attacker = players[turns%2]
		damage := rand.Intn(21) + 10 // 10~30
		players[(turns+1)%2].Hp -= damage
		turns += 1
	}

	return &BattleResult{
		MatchId: matchId,
		Player1: players[0].UserId,
		Player2: players[1].UserId,
		Winner: attacker.UserId,
		Turns: turns+1,
	}
}

func (s *MatchServer) GetBattleResult(ctx context.Context, req *pb.BattleResultRequest) (*pb.BattleResultResponse, error) {
	battleResult, ok := s.matchmaker.battleResults[req.MatchId]
	if !ok {
		log.Printf("Battle result not found for matchId: %s", req.MatchId)
		return nil, status.Errorf(codes.NotFound, "Battle result not found for matchId: %s", req.MatchId)
	}
	return &pb.BattleResultResponse{
		MatchId: battleResult.MatchId,
		Player1: battleResult.Player1,
		Player2: battleResult.Player2,
		Winner: battleResult.Winner,
		Turns: int32(battleResult.Turns),
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
	pb.RegisterMatchmakingServiceServer(grpcServer, &MatchServer{
		matchmaker: matchmaker,
	})

	log.Println("gRPC server is listening on port 50051...")
	grpcServer.Serve(lis)
}