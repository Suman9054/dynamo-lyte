package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/suman9054/zeone/pkg/proto"

)

type server struct {
	pb.UnimplementedQueryServiceServer
}


func (s *server) ExecuteQuery(
	ctx context.Context,
	req *pb.QueryRequest,
) (*pb.QueryResponse, error) {

	log.Println("Received query:", req.Query)
	log.Println("Parameters:", req.Parameters)

	
	return &pb.QueryResponse{
		Status: pb.Status_SUCCESS,
	}, nil
}

func main() {
	
	lis, err := net.Listen("tcp", "127.0.0.1:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	
	grpcServer := grpc.NewServer()

	
	pb.RegisterQueryServiceServer(grpcServer, &server{})

	log.Println("🚀 gRPC server running on 127.0.0.1:50051")

	
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
