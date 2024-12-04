package store

import (
	"cloudbee/internal/model"
	"fmt"
	"sync"
)

type BookingStore struct {
	mu      sync.RWMutex
	tickets map[string]model.Ticket           // email -> ticket
	seats   map[string]map[int32]string 	// section -> seat number -> email
}

func NewBookingStore() *BookingStore {
	return &BookingStore{
		tickets: make(map[string]model.Ticket),
		seats: map[string]map[int32]string{
			model.SectionA: make(map[int32]string),
			model.SectionB: make(map[int32]string),
		},
	}
}

func (bs *BookingStore) PurchaseTicket(user model.User) (*model.Ticket, error) {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	if _, exists := bs.tickets[user.Email]; exists {
		return nil, fmt.Errorf("user already has a ticket")
	}

	// Find available seat
	seat, err := bs.findAvailableSeat()
	if err != nil {
		return nil, err
	}

	ticket := model.Ticket{
		From:  model.FromStation,
		To:    model.ToStation,
		User:  user,
		Price: model.TicketPrice,
		Seat:  *seat,
	}

	// Store ticket and seat allocation
	bs.tickets[user.Email] = ticket
	bs.seats[seat.Section][seat.Number] = user.Email

	return &ticket, nil
}

func (bs *BookingStore) GetReceipt(email string) (*model.Ticket, error) {
	bs.mu.RLock()
	defer bs.mu.RUnlock()

	ticket, exists := bs.tickets[email]
	if !exists {
		return nil, fmt.Errorf("no ticket found for email: %s", email)
	}

	return &ticket, nil
}

func (bs *BookingStore) GetSectionUsers(section string) ([]model.Ticket, error) {
	bs.mu.RLock()
	defer bs.mu.RUnlock()

	if section != model.SectionA && section != model.SectionB {
		return nil, fmt.Errorf("invalid section: %s", section)
	}

	var tickets []model.Ticket
	for _, ticket := range bs.tickets {
		if ticket.Seat.Section == section {
			tickets = append(tickets, ticket)
		}
	}

	return tickets, nil
}

func (bs *BookingStore) RemoveUser(email string) error {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	ticket, exists := bs.tickets[email]
	if !exists {
		return fmt.Errorf("no ticket found for email: %s", email)
	}

	delete(bs.seats[ticket.Seat.Section], ticket.Seat.Number)
	delete(bs.tickets, email)

	return nil
}

func (bs *BookingStore) ModifySeat(email string, newSeat model.Seat) (*model.Ticket, error) {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	if newSeat.Section != model.SectionA && newSeat.Section != model.SectionB {
		return nil, fmt.Errorf("invalid section: %s", newSeat.Section)
	}

	// Check if seat is available
	if occupantEmail, occupied := bs.seats[newSeat.Section][newSeat.Number]; occupied {
		if occupantEmail != email {
			return nil, fmt.Errorf("seat %d in section %s is already occupied", newSeat.Number, newSeat.Section)
		}
	}

	ticket, exists := bs.tickets[email]
	if !exists {
		return nil, fmt.Errorf("no ticket found for email: %s", email)
	}

	delete(bs.seats[ticket.Seat.Section], ticket.Seat.Number)

	// Update seat
	ticket.Seat = newSeat
	bs.tickets[email] = ticket
	bs.seats[newSeat.Section][newSeat.Number] = email

	return &ticket, nil
}

func (bs *BookingStore) findAvailableSeat() (*model.Seat, error) {
	sections := []string{model.SectionA, model.SectionB}
	for _, section := range sections {
		for seatNum := int32(1); seatNum <= model.SeatsPerSection; seatNum++ {
			if _, occupied := bs.seats[section][seatNum]; !occupied {
				return &model.Seat{Section: section, Number: seatNum}, nil
			}
		}
	}
	return nil, fmt.Errorf("no available seats")
}

func (bs *BookingStore) GetStat() (*model.StatResponse, error) {
	bs.mu.RLock()
	defer bs.mu.RUnlock()

	var user []model.User
	for _, ticket := range bs.tickets {
		user = append(user, ticket.User)
	}

	var seat []model.Seat
	for id, section := range bs.seats {
		for seatNum, _ := range section {
			seat = append(seat, model.Seat{Section: id, Number: seatNum})
		}
	}

	var tickets []model.Ticket
	for _, ticket := range bs.tickets {
		tickets = append(tickets, ticket)
	}

	return &model.StatResponse{
		User: user,
		Seat: seat,
		Ticket: tickets,
	}, nil
}
