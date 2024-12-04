package main

import (
	"fmt"
	"log"
	"net"

	istore "cloudbee/internal/store"
	"cloudbee/internal/service"
	pb "cloudbee/proto/booking/v1"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	svc := service.NewBookingServer(istore.NewBookingStore())
	pb.RegisterBookingServiceServer(s, svc)

	fmt.Printf("Server listening on port %s\n", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
