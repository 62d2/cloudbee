package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "cloudbee/proto/booking/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	address = "localhost:50051"
)

func stat(ctx context.Context, c pb.BookingServiceClient) {
	fmt.Printf("-------------------------------------------------\n")
	fmt.Println("\n >>>> Current Data of Database : ")
	response, err := c.GetStat(ctx, &pb.GetStatRequest{})
	if err != nil {
		log.Fatalf("could not get Database : %v", err)
	}
	fmt.Printf("Ticket stat: %+v\n", response)
	fmt.Printf("-------------------------------------------------\n")
}

func main() {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewBookingServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// UseCase 1: Purchase a ticket
	fmt.Println("\n1. Purchasing a ticket...")
	purchaseResp, err := c.PurchaseTicket(ctx, &pb.PurchaseTicketRequest{
		User: &pb.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
		},
	})
	if err != nil {
		log.Fatalf("could not purchase ticket: %v", err)
	}
	fmt.Printf("Ticket purchased: %+v\n", purchaseResp.Ticket)
	stat(ctx, c)

	// UseCase 2: Get receipt
	fmt.Println("\n2. Getting receipt...")
	receiptResp, err := c.GetReceipt(ctx, &pb.GetReceiptRequest{
		Email: "john.doe@example.com",
	})
	if err != nil {
		log.Fatalf("could not get receipt: %v", err)
	}
	fmt.Printf("Receipt details: %+v\n", receiptResp.Ticket)
	stat(ctx, c)

	// UseCase 3: Get users in section A
	fmt.Println("\n3. Getting users in section A...")
	sectionResp, err := c.GetSectionUsers(ctx, &pb.GetSectionUsersRequest{
		Section: "A",
	})
	if err != nil {
		log.Fatalf("could not get section users: %v", err)
	}
	fmt.Printf("Users in section A: %+v\n", sectionResp.UserSeats)
	stat(ctx, c)

	// UseCase 4: Modify seat
	fmt.Println("\n4. Modifying seat...")
	modifySeatResp, err := c.ModifySeat(ctx, &pb.ModifySeatRequest{
		Email: "john.doe@example.com",
		NewSeat: &pb.Seat{
			Section: "B",
			Number:  1,
		},
	})
	if err != nil {
		log.Fatalf("could not modify seat: %v", err)
	}
	fmt.Printf("Updated ticket after seat modification: %+v\n", modifySeatResp.UpdatedTicket)
	stat(ctx, c)

	// UseCase 5: Remove user
	fmt.Println("\n5. Removing user...")
	removeResp, err := c.RemoveUser(ctx, &pb.RemoveUserRequest{
		Email: "john.doe@example.com",
	})
	if err != nil {
		log.Fatalf("could not remove user: %v", err)
	}
	fmt.Printf("User removed successfully: %v\n", removeResp.Success)
	stat(ctx, c)
}
