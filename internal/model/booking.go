package model


const (
	SectionA 		= "A"
	SectionB 		= "B"
	FromStation		= "London"
	ToStation		= "France"
	SeatsPerSection = 50
	TicketPrice 	= 20.0
)

type User struct {
	FirstName string
	LastName  string
	Email     string
}

type Seat struct {
	Section string
	Number  int32
}

type Ticket struct {
	From  string
	To    string
	User  User
	Price float64
	Seat  Seat
}

type StatResponse struct {
	User []User
	Seat []Seat
	Ticket []Ticket
}
