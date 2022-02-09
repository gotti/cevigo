package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	pb "github.com/gotti/cevigo/spec"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	addr := flag.String("addr", "localhost:10001", "server address")
	flag.Parse()

	fmt.Printf("connecting to %s\n", *addr)
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v\n", err)
	}
	defer conn.Close()

	c := pb.NewTtsClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	r, err := c.CreateWav(ctx, &pb.CevioTtsRequest{
		Text:   "電気通信大学は、武蔵野の緑溢れる東京都調布市にある国立大学です。",
		Cast:   "さとうささら",
		Volume: 50,
		Speed:  50,
		Pitch:  50,
		Alpha:  50,
		Intro:  50,
		Emotions: map[string]uint32{
			"普通":  100,
			"元気":  0,
			"怒り":  0,
			"哀しみ": 0,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	os.WriteFile("./out.wav", r.Audio, os.ModePerm)
}
