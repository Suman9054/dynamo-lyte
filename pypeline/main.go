package main

import (
	"context"
	"fmt"

	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/structpb"

	pb "github.com/suman9054/zeone/pkg/proto"
	"github.com/suman9054/zeone/pkg/store"
)

type server struct {
	pb.UnimplementedQueryServiceServer
}

var Questore = store.Newstore()

func (s *server) PutQuery(
	ctx context.Context,
	req *pb.QueryRequest,
) (*pb.QueryResponse, error) {

	if req.Query == "" || req.Parameters == nil {
		return &pb.QueryResponse{
			Status: pb.Status_STATUS_UNSPECIFIED,
		}, nil
	}
	// put your query execution logic here
	result, err := structpb.NewStruct(req.Parameters.AsMap())

	if err != nil {
		return &pb.QueryResponse{
			Status: pb.Status_FAILURE,
		}, nil
	}
	Questore.Putdata(req.Query, result)

	return &pb.QueryResponse{
		Status: pb.Status_SUCCESS,
	}, nil
}

func (s *server) GetQuery(
	ctx context.Context,
	req *pb.GetQueryRequest,
) (*pb.GetQueryResponse, error) {
	if req.Querykey == "" {
		return &pb.GetQueryResponse{
			Status: pb.Status_STATUS_UNSPECIFIED,
		}, nil
	}

	value, ok := Questore.Getdata(req.Querykey)
	if !ok {
		return &pb.GetQueryResponse{
			Status: pb.Status_NOT_FOUND,
		}, nil
	}

	return &pb.GetQueryResponse{
		Status:      pb.Status_SUCCESS,
		Queryresult: value,
	}, nil
}

func (s *server) SubscribeToQuery(
	req *pb.GetQueryRequest,
	stream pb.QueryService_SubscribeToQueryServer,
) error {
	key := req.Querykey
	ch, clean := Questore.Subscribequery(key)
	defer clean()

	for {
		select {
		case <-stream.Context().Done():
			// Client disconnected
			fmt.Println("Client disconnected")
			return stream.Context().Err()

		case update, ok := <-ch:
			if !ok {
				// Channel closed by server
				fmt.Println("Query subscription channel closed")
				return nil
			}

			result := &pb.GetQueryResponse{
				Queryresult: update,
			}

			if err := stream.Send(result); err != nil {
				fmt.Println("Send error:", err)
				return err
			}
		}
	}
}

func main() {

	lis, err := net.Listen("tcp", "127.0.0.1:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	go Questore.Expirewatcher()

	grpcServer := grpc.NewServer()

	pb.RegisterQueryServiceServer(grpcServer, &server{})

	log.Println("🚀 gRPC server running on 127.0.0.1:50051")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
