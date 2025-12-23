package main

import (
	"context"
	"log"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/suman9054/zeone/pkg/proto"
)

type server struct {
	pb.UnimplementedQueryServiceServer
}

var( 
	Questore map[string]map[string]string
	mu       sync.Mutex
)

func (s *server) ExecuteQuery(
	ctx context.Context,
	req *pb.QueryRequest,
) (*pb.QueryResponse, error) {
  
 if req.Query == "" || req.Parameters == nil {
		return &pb.QueryResponse{
			Status: pb.Status_STATUS_UNSPECIFIED,
		}, nil
	}

	mu.Lock()
	Questore[req.Query] = req.Parameters
	mu.Unlock()

	return &pb.QueryResponse{
		Status: pb.Status_SUCCESS,
	}, nil
}

func (s *server) GetQuery(
	ctx context.Context,
	req *pb.Parameters,
) (*pb.Queryresponse, error){
	
	if req.Query == "" {
		return &pb.Queryresponse{
			Status: pb.Status_STATUS_UNSPECIFIED,
		}, nil
	}
	
	mu.Lock()
	response,ok := Questore[req.Query]
	mu.Unlock()
	if !ok {
		return &pb.Queryresponse{
			Status: pb.Status_STATUS_UNSPECIFIED,
		}, nil
	}

	return &pb.Queryresponse{
		Response: response,
		Status:   pb.Status_SUCCESS,
	}, nil

}

func (s *server) ConsumeQuery(
	req *pb.Parameters,
	stream grpc.ServerStreamingServer[pb.Queryresponse],
) error {

	if req.Query == "" {
		return status.Error(codes.InvalidArgument, "query is empty")
	}

	for {
		// stop if client disconnects
		if stream.Context().Err() != nil {
			return stream.Context().Err()
		}

		mu.Lock()
		response, ok := Questore[req.Query]
		mu.Unlock()

		if ok {
			err := stream.Send(&pb.Queryresponse{
				Response: response,
				Status:   pb.Status_SUCCESS,
			})
			if err != nil {
				return err
			}
		}

		// prevent busy loop
		time.Sleep(500 * time.Millisecond)
	}
}

func main() {
	Questore = make(map[string]map[string]string)
	
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
