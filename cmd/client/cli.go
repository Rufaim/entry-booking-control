package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	pb "github.com/Rufaim/entry_booking_control/cmd/message"
)

func RunCLIApplication(client pb.LabVisitsServiceClient) {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	switch os.Args[1] {
	case HelpCommand:
		printUsage()
	case BookCommand:
		BookingHandler(client)
	case VisitsCommand:
		VisitsHandler(client)
	default:
		printUsage()
	}
}

func BookingHandler(client pb.LabVisitsServiceClient) {
	flagset := flag.NewFlagSet(BookCommand, flag.PanicOnError)

	weekdayStr := flagset.String(BookCommandSetWeekdayFlag, "", "")
	del := flagset.Bool(BookCommandDeleteFlag, false, "")

	if err := flagset.Parse(os.Args[2:]); err != nil {
		printUsage()
		panic(err)
	}

	weekday, ok := pb.Weekday_value[*weekdayStr]
	if !ok {
		printUsage()
		return
	}

	ctx := context.Background()
	timestamp := time.Now().Format(TimeFormat)
	usr := GetUserIdentifier()
	uvs := &pb.UserVisitSet{
		User: &pb.User{
			Id: usr,
		},
		Day:       pb.Weekday(weekday),
		Timestamp: timestamp,
	}
	var (
		result *pb.UserVisitSetResult
		err    error
	)
	if *del {
		result, err = client.DelUserVisit(ctx, uvs)

	} else {
		result, err = client.SetUserVisit(ctx, uvs)
	}

	PanicOnError(err)
	if result.Status != pb.UserVisitSetResult_OK {
		fmt.Println(result.Text)
	}
}

func VisitsHandler(client pb.LabVisitsServiceClient) {
	flagset := flag.NewFlagSet(VisitsCommand, flag.PanicOnError)

	isTotal := flagset.Bool(VisitsCommandAllFlag, false, "")

	if err := flagset.Parse(os.Args[2:]); err != nil {
		printUsage()
		panic(err)
	}

	ctx := context.Background()
	if *isTotal {
		allStats, err := client.GetAllVisitsReport(ctx, &pb.Void{})
		PanicOnError(err)
		PrintTotalReport(allStats)
		return
	}

	usr := GetUserIdentifier()
	usrStat, err := client.GetVisits(ctx, &pb.User{
		Id: usr,
	})
	PanicOnError(err)
	PrintUserReport(usrStat)
}
