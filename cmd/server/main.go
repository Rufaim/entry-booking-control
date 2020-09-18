package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	pb "github.com/Rufaim/entry_booking_control/cmd/message"
	"google.golang.org/grpc"
)

func main() {
	// Logging initialization
	file, err := os.OpenFile(LogFilename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	logger := log.New(file, "", log.LstdFlags|log.Lshortfile)

	lis, err := net.Listen("tcp", Port)
	if err != nil {
		panic(fmt.Errorf("failed to listen: %v", err))
	}
	logger.Printf("Listening to port %s\n", Port)
	s := grpc.NewServer()
	server, err := NewServer(logger)
	if err != nil {
		panic(fmt.Errorf("failed to initialize server: %v", err))
	}

	pb.RegisterLabVisitsServiceServer(s, server)
	tickerSync := time.Tick(DatabaseSyncTime)
	tickerPrune := time.Tick(HistoryPrunningTime)
	go func() {
		for {
			<-tickerSync
			if err := server.Sync(); err != nil {
				panic(fmt.Errorf("failed to syncronize server: %v", err))
			}
		}
	}()
	go func() {
		for {
			<-tickerPrune
			if err := server.PruneHistory(); err != nil {
				panic(fmt.Errorf("failed to prune history: %v", err))
			}
		}
	}()

	if err := s.Serve(lis); err != nil {
		panic(fmt.Errorf("failed to serve: %v", err))
	}
}
