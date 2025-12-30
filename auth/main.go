package main

import (
	"context"
	"fmt"

	pb "gihub.com/suman9054/dynamo-lyte/auth/pkg/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("Failed to connect:", err)
	}
	defer conn.Close()
	// Further client logic goes here
	//....
	client := pb.NewQueryServiceClient(conn)
	ctx, cncel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cncel()

}
