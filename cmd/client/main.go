package main

import (
	pb "github.com/Rufaim/entry_booking_control/cmd/message"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial(Address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := pb.NewLabVisitsServiceClient(conn)

	RunCLIApplication(client)
}
