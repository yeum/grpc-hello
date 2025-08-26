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

type SessionServer struct {
	pb.UnimplementedSessionServiceServer
	sm *SessionManager
}

type SessionManager struct {
	players map[string]chan *Player
}

func NewSessionManager {
	return &SessionManager{
		players: make(map[string]chan *Player)
	}
}

type Player struct {
	sessionId string
	userId string
	expiresAt int
	locX int
	locT int
	charNum int
}

func NewPlayer(userId strig) {
	return &NewPlayer{
		sessionId: uuid.New.String(),
		userId: userId,
		expiresAt: time.Now().Add(time.Hour).Unix(),
		locX: 0,
		locY: 0,
		charNum: rand.Intn(15)+1
	}
}

func (s* SessionServer) CreateSession(ctx context.Context, req *pb.CreateSessionRequest) (*pb.CreateSessionResponse, error) {
	newPlayer = s.sm.AddPlayer(req.UserId)
	return &pb.CreateSessionResponse{
		SessionId: newPlayer.sessionId,
		CharNum: newPlayer.charNum,
		ExpiresAt: newPlayer.expiresAt,
		LocationX: newPlayer.locX,
		LocationY: newPlayer.locY
	}, nil
}

func (sm* SessionManager) AddPlayer(userId string) *Player {
	newPlayer := NewPlayer(userId)
	sm.players[newPlayer.sessionId] = newPlayer
	return newPlayer
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	sessionManager := NewSessionManager()
	pb.RegisterSessionServiceServer(grpcServer, &sessionManager{})

	log.Println("gRPC server is listening on port 50051...")
	grpcServer.Serve(lis)
}