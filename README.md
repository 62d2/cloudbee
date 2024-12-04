# Train Booking System

A gRPC-based train booking system that allows users to book tickets for trains between London and France.

## Features

- Purchase train tickets
- View ticket receipts
- View seat allocations by section
- Remove users from train
- Modify user's seat

## Prerequisites

- Go 1.23.4 or later
- Protocol Buffers compiler
- gRPC tools

## API Endpoints

The following gRPC endpoints are available:

1. `PurchaseTicket` - Purchase a train ticket
2. `GetReceipt` - Get ticket receipt details
3. `GetSectionUsers` - View users and seats by section
4. `RemoveUser` - Remove a user from the train
5. `ModifySeat` - Modify a user's seat

## Implementation Details

- The system uses in-memory storage for data persistence
- Train has two sections: A and B
- Fixed ticket price of $20
- Basic validation for email and seat allocation


## How to run the project

- Install Go 1.23.4 or later
- Make sure you have all the protobuf in you PC
  - ``` brew install protobuf ```
- Generate the gRPC code
```bash
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/booking/v1/booking.proto
```
- Install the dependencies
```bash
go mod download && go mod tidy
```
- Run the server
```bash
go run cmd/server/main.go
```
- Run the client
```bash
go run cmd/client/main.go
```
# cloudbee
