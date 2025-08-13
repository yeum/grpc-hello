package main

import (
	"context"
	"log"
	"time"
	"os"
	"fmt"

	pb "grpc-hello/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
    "google.golang.org/grpc/codes"
)

func main() {
	addr := os.Getenv("SERVER_ADDR")
	if addr == "" {
		addr = "localhost:50051"
	}
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	helloClient := pb.NewHelloServiceClient(conn)
	matchClient := pb.NewMatchmakingServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userId := os.Args[1]
	matchId := ""

	{
		res, err := helloClient.SayHello(ctx, &pb.HelloRequest{Name: userId})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}

		log.Printf("Greeting: %s", res.GetMessage())
	}
	{
		joinRes, err := matchClient.JoinQueue(ctx, &pb.JoinRequest{UserId: userId})
		if err != nil {
			log.Fatalf("could not join: %v", err)
		}
		if joinRes.GetIsMatched() {
			matchId = joinRes.GetMatchId()
			log.Printf("Match success! MatchId: %s, OpponentId: %s", matchId, joinRes.GetOpponentId())
			
			time.Sleep(time.Second)
			getBattleRes, err := matchClient.GetBattleResult(ctx, &pb.BattleResultRequest{MatchId: matchId})
			if err != nil {
				stat, _ := status.FromError(err)
				if stat.Code() == codes.NotFound {
					fmt.Println("해당 매치를 찾을 수 없습니다.")
					return
				} else {
					log.Fatalf("server error: %v", err)
				}
			}
			log.Printf("[%s] Player1: %s, Player2: %s, Winner: %s, Turns: %v", getBattleRes.GetMatchId(), getBattleRes.GetPlayer1(), getBattleRes.GetPlayer2(), getBattleRes.GetWinner(), getBattleRes.GetTurns())
		} else {
			log.Printf("Waiting to match...")
		}
	}
}