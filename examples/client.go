package main

import (
	"context"
	"io"
	"log"

	"github.com/Andrewmakmaer/SystemInfoFluxDaemon/internal/server/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient(":8765", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewDaemonClient(conn)

	ctx, close := context.WithCancel(context.Background())
	defer close()

	req := pb.StreamRequest{StatsType: "cpu", SecondsRange: 10, SecondsDelay: 5}
	stream, err := client.EnableStatStream(ctx, &req)
	if err != nil {
		log.Fatal(err)
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Fatal(err)
		}
		stat := resp.GetCpuStats()

		log.Printf("CPU stats | sys: %f | usr: %f | idle: %f | iowait: %f\n", stat.Sys, stat.Usr, stat.Idle, stat.Iowait)
	}
}
