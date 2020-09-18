package main

import (
	"fmt"
	"os"

	pb "github.com/Rufaim/entry_booking_control/cmd/message"
)

func printUsage() {
	fmt.Println("Usage:")
	fmt.Printf("%s %s -%s <weekday> -- books a slot at given weekday (applicable weekdays are %s, %s, %s, %s, %s)\n",
		os.Args[0], BookCommand, BookCommandSetWeekdayFlag, pb.Weekday_Mon.String(),
		pb.Weekday_Tue.String(), pb.Weekday_Wed.String(),
		pb.Weekday_Thu.String(), pb.Weekday_Fri.String())
	fmt.Printf("%s %s -%s -%s <weekday> -- deletes booking of a slot at given weekday (applicable weekdays are %s, %s, %s, %s, %s)\n",
		os.Args[0], BookCommand, BookCommandDeleteFlag, BookCommandSetWeekdayFlag,
		pb.Weekday_Mon.String(), pb.Weekday_Tue.String(), pb.Weekday_Wed.String(),
		pb.Weekday_Thu.String(), pb.Weekday_Fri.String())
	fmt.Printf("%s %s -- prints current user's bookings guring a week period\n", os.Args[0], VisitsCommand)
	fmt.Printf("%s %s -%s -- prints info of all bookings guring a week period\n", os.Args[0], VisitsCommand, VisitsCommandAllFlag)
	fmt.Printf("%s [%s] -- prints this help\n", os.Args[0], HelpCommand)
}

func PrintTotalReport(allStats *pb.AllUsersStat) {
	days := make([][]string, len(pb.Weekday_name))
	for i := range days {
		days[i] = make([]string, 0)
	}
	for _, stat := range allStats.Stat {
		for _, visit := range stat.Visits {
			days[int(visit.Day)] = append(days[int(visit.Day)], stat.User.GetId())
		}
	}
	for i, d := range days {
		if len(d) == 0 {
			continue
		}
		fmt.Printf("%s:\n", pb.Weekday_name[int32(i)])
		for _, id := range d {
			fmt.Printf("  %s\n", id)
		}
	}
	fmt.Println("\nBooking info:")
	for i := range days {
		day := pb.Weekday_name[int32(i)]
		val := allStats.BookingAmount[day]
		fmt.Printf("%s:%d   ", day, val)
	}
	fmt.Println("")
}

func PrintUserReport(usrStat *pb.UserStat) {
	fmt.Printf("%s:\n", usrStat.User.GetId())
	if len(usrStat.Visits) == 0 {
		fmt.Println("<None>")
		return
	}
	for _, visit := range usrStat.Visits {
		fmt.Printf("%s   ", visit.Day.String())
	}
	fmt.Println("")
}
