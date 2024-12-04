package service

import (
	"context"
	"testing"

	istore "cloudbee/internal/store"
	pb "cloudbee/proto/booking/v1"
)

var (
	store = istore.NewBookingStore()
)

func TestNewBookingServer(t *testing.T) {
	server := NewBookingServer(store)
	if server == nil {
		t.Fatal("Expected non-nil BookingServer")
	}
	if server.store == nil {
		t.Error("Expected non-nil BookingStore in server")
	}
}

func TestPurchaseTicketService(t *testing.T) {
	server := NewBookingServer(store)
	ctx := context.Background()

	tests := []struct {
		name    string
		req     *pb.PurchaseTicketRequest
		wantErr bool
	}{
		{
			name: "Valid purchase",
			req: &pb.PurchaseTicketRequest{
				User: &pb.User{
					FirstName: "John",
					LastName:  "Doe",
					Email:     "john@example.com",
				},
			},
			wantErr: false,
		},
		{
			name: "Duplicate purchase",
			req: &pb.PurchaseTicketRequest{
				User: &pb.User{
					FirstName: "John",
					LastName:  "Doe",
					Email:     "john@example.com",
				},
			},
			wantErr: true,
		},
		{
			name: "Missing email",
			req: &pb.PurchaseTicketRequest{
				User: &pb.User{
					FirstName: "John",
					LastName:  "Doe",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := server.PurchaseTicket(ctx, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("PurchaseTicket() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if resp == nil {
					t.Error("Expected non-nil response")
					return
				}
				if resp.Ticket == nil {
					t.Error("Expected non-nil ticket in response")
					return
				}
				if resp.Ticket.User.Email != tt.req.User.Email {
					t.Error("Response ticket email does not match request")
				}
			}
		})
	}
}

func TestGetReceiptService(t *testing.T) {
	server := NewBookingServer(store)
	ctx := context.Background()

	// First purchase a ticket
	purchaseReq := &pb.PurchaseTicketRequest{
		User: &pb.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john@example.com",
		},
	}
	purchaseResp, _ := server.PurchaseTicket(ctx, purchaseReq)

	tests := []struct {
		name    string
		req     *pb.GetReceiptRequest
		wantErr bool
	}{
		{
			name: "Valid receipt request",
			req: &pb.GetReceiptRequest{
				Email: "john@example.com",
			},
			wantErr: false,
		},
		{
			name: "Non-existent user",
			req: &pb.GetReceiptRequest{
				Email: "nonexistent@example.com",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := server.GetReceipt(ctx, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetReceipt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if resp == nil {
					t.Error("Expected non-nil response")
					return
				}
				if resp.Ticket == nil {
					t.Error("Expected non-nil ticket in response")
					return
				}
				if resp.Ticket.User.Email != purchaseResp.Ticket.User.Email {
					t.Error("Receipt ticket does not match purchased ticket")
				}
			}
		})
	}
}

func TestGetSectionUsersService(t *testing.T) {
	server := NewBookingServer(store)
	ctx := context.Background()

	// Add some users
	users := []*pb.User{
		{FirstName: "John", LastName: "Doe", Email: "john@example.com"},
		{FirstName: "Jane", LastName: "Doe", Email: "jane@example.com"},
	}

	for _, user := range users {
		_, _ = server.PurchaseTicket(ctx, &pb.PurchaseTicketRequest{User: user})
	}

	tests := []struct {
		name    string
		req     *pb.GetSectionUsersRequest
		wantErr bool
	}{
		{
			name: "Valid section A",
			req: &pb.GetSectionUsersRequest{
				Section: "A",
			},
			wantErr: false,
		},
		{
			name: "Valid section B",
			req: &pb.GetSectionUsersRequest{
				Section: "B",
			},
			wantErr: false,
		},
		{
			name: "Invalid section",
			req: &pb.GetSectionUsersRequest{
				Section: "C",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := server.GetSectionUsers(ctx, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSectionUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if resp == nil {
					t.Error("Expected non-nil response")
					return
				}
				if resp.UserSeats == nil {
					t.Error("Expected non-nil user seats in response")
				}
			}
		})
	}
}

func TestRemoveUserService(t *testing.T) {
	server := NewBookingServer(store)
	ctx := context.Background()

	// First purchase a ticket
	purchaseReq := &pb.PurchaseTicketRequest{
		User: &pb.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john@example.com",
		},
	}
	_, _ = server.PurchaseTicket(ctx, purchaseReq)

	tests := []struct {
		name    string
		req     *pb.RemoveUserRequest
		wantErr bool
	}{
		{
			name: "Valid removal",
			req: &pb.RemoveUserRequest{
				Email: "john@example.com",
			},
			wantErr: false,
		},
		{
			name: "Non-existent user",
			req: &pb.RemoveUserRequest{
				Email: "nonexistent@example.com",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := server.RemoveUser(ctx, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("RemoveUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if resp == nil {
					t.Error("Expected non-nil response")
					return
				}
				if !resp.Success {
					t.Error("Expected success to be true")
				}
			}
		})
	}
}

func TestModifySeatService(t *testing.T) {
	server := NewBookingServer(store)
	ctx := context.Background()

	// First purchase a ticket
	purchaseReq := &pb.PurchaseTicketRequest{
		User: &pb.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john@example.com",
		},
	}
	purchaseResp, _ := server.PurchaseTicket(ctx, purchaseReq)

	tests := []struct {
		name    string
		req     *pb.ModifySeatRequest
		wantErr bool
	}{
		{
			name: "Valid modification",
			req: &pb.ModifySeatRequest{
				Email: "john@example.com",
				NewSeat: &pb.Seat{
					Section: "B",
					Number:  1,
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid section",
			req: &pb.ModifySeatRequest{
				Email: "john@example.com",
				NewSeat: &pb.Seat{
					Section: "C",
					Number:  1,
				},
			},
			wantErr: true,
		},
		{
			name: "Non-existent user",
			req: &pb.ModifySeatRequest{
				Email: "nonexistent@example.com",
				NewSeat: &pb.Seat{
					Section: "A",
					Number:  1,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := server.ModifySeat(ctx, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ModifySeat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if resp == nil {
					t.Error("Expected non-nil response")
					return
				}
				if resp.UpdatedTicket == nil {
					t.Error("Expected non-nil updated ticket in response")
					return
				}
				if resp.UpdatedTicket.User.Email != purchaseResp.Ticket.User.Email {
					t.Error("Updated ticket user does not match original user")
				}
				if resp.UpdatedTicket.Seat.Section != tt.req.NewSeat.Section ||
					resp.UpdatedTicket.Seat.Number != tt.req.NewSeat.Number {
					t.Error("Seat was not properly modified")
				}
			}
		})
	}
}
