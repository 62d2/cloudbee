package service

import (
	"context"

	"cloudbee/internal/model"
	istore "cloudbee/internal/store"
	pb "cloudbee/proto/booking/v1"
)

type BookingServer struct {
	pb.UnimplementedBookingServiceServer
	store *istore.BookingStore
}

func NewBookingServer(store *istore.BookingStore) *BookingServer {
	return &BookingServer{
		store: istore.NewBookingStore(),
	}
}

func (s *BookingServer) PurchaseTicket(ctx context.Context, req *pb.PurchaseTicketRequest) (*pb.PurchaseTicketResponse, error) {
	user := model.User{
		FirstName: req.User.FirstName,
		LastName:  req.User.LastName,
		Email:     req.User.Email,
	}

	ticket, err := s.store.PurchaseTicket(user)
	if err != nil {
		return nil, err
	}

	return &pb.PurchaseTicketResponse{
		Ticket: convertTicketToProto(ticket),
	}, nil
}

func (s *BookingServer) GetReceipt(ctx context.Context, req *pb.GetReceiptRequest) (*pb.GetReceiptResponse, error) {
	ticket, err := s.store.GetReceipt(req.Email)
	if err != nil {
		return nil, err
	}

	return &pb.GetReceiptResponse{
		Ticket: convertTicketToProto(ticket),
	}, nil
}

func (s *BookingServer) GetSectionUsers(ctx context.Context, req *pb.GetSectionUsersRequest) (*pb.GetSectionUsersResponse, error) {
	tickets, err := s.store.GetSectionUsers(req.Section)
	if err != nil {
		return nil, err
	}

	userSeats := make([]*pb.GetSectionUsersResponse_UserSeat, 0, len(tickets))
	for _, ticket := range tickets {
		userSeats = append(userSeats, &pb.GetSectionUsersResponse_UserSeat{
			User: &pb.User{
				FirstName: ticket.User.FirstName,
				LastName:  ticket.User.LastName,
				Email:     ticket.User.Email,
			},
			Seat: &pb.Seat{
				Section: ticket.Seat.Section,
				Number:  ticket.Seat.Number,
			},
		})
	}

	return &pb.GetSectionUsersResponse{
		UserSeats: userSeats,
	}, nil
}

func (s *BookingServer) RemoveUser(ctx context.Context, req *pb.RemoveUserRequest) (*pb.RemoveUserResponse, error) {
	err := s.store.RemoveUser(req.Email)
	if err != nil {
		return &pb.RemoveUserResponse{Success: false}, err
	}

	return &pb.RemoveUserResponse{Success: true}, nil
}

func (s *BookingServer) ModifySeat(ctx context.Context, req *pb.ModifySeatRequest) (*pb.ModifySeatResponse, error) {
	newSeat := model.Seat{
		Section: req.NewSeat.Section,
		Number:  req.NewSeat.Number,
	}

	ticket, err := s.store.ModifySeat(req.Email, newSeat)
	if err != nil {
		return nil, err
	}

	return &pb.ModifySeatResponse{
		UpdatedTicket: convertTicketToProto(ticket),
	}, nil
}

func (s *BookingServer) GetStat(ctx context.Context, req *pb.GetStatRequest) (*pb.GetStatResponse, error) {

	response, err := s.store.GetStat()
	if err != nil {
		return nil, err
	}

	seats := make([]*pb.Seat, 0, len(response.Seat))
	for _, seat := range response.Seat {
		seats = append(seats, &pb.Seat{
			Section: seat.Section,
			Number:  seat.Number,
		})
	}
	users := make([]*pb.User, 0, len(response.User))
	for _, user := range response.User {
		users = append(users, &pb.User{
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
		})
	}
	tickets := make([]*pb.Ticket, 0, len(response.Ticket))
	for _, ticket := range response.Ticket {
		tickets = append(tickets, &pb.Ticket{
			From:     ticket.From,
			To:       ticket.To,
		})
	}



	return &pb.GetStatResponse{
		User:   users,
		Seat:   seats,
		Ticket: tickets,
	}, nil
}

func convertTicketToProto(ticket *model.Ticket) *pb.Ticket {
	return &pb.Ticket{
		From: ticket.From,
		To:   ticket.To,
		User: &pb.User{
			FirstName: ticket.User.FirstName,
			LastName:  ticket.User.LastName,
			Email:     ticket.User.Email,
		},
		Price: ticket.Price,
		Seat: &pb.Seat{
			Section: ticket.Seat.Section,
			Number:  ticket.Seat.Number,
		},
	}
}
