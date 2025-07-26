package main

import (
	"context"
	"log"
	"time"
	"os"

	pb "grpc-hello/proto"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	helloClient := pb.NewHelloServiceClient(conn)
	matchClient := pb.NewMatchmakingServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userId := os.Args[1]

	{
		res, err := helloClient.SayHello(ctx, &pb.HelloRequest{Name: userId})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}

		log.Printf("Greeting: %s", res.GetMessage())
	}
	{
		res, err := matchClient.JoinQueue(ctx, &pb.JoinRequest{UserId: userId})
		if err != nil {
			log.Fatalf("could not join: %v", err)
		}
		if res.GetIsMatched() {
			log.Printf("Match success! MatchId: %s, OpponentId: %s", res.GetMatchId(), res.GetOpponentId())
		} else {
			log.Printf("Waiting to match...")
		}
	}
	
}