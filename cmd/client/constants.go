package main

import "time"

const Address = "localhost:6061"

const TimeFormat = time.RFC3339

// help command
const (
	HelpCommand = "help"
)

// book command
const (
	BookCommand               = "book"
	BookCommandSetWeekdayFlag = "weekday"
	BookCommandDeleteFlag     = "d"
)

// visits command
const (
	VisitsCommand        = "visits"
	VisitsCommandAllFlag = "a"
)
