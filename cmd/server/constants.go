package main

import (
	"time"

	"github.com/pkg/errors"
)

const (
	DatabaseFilename = "base.db"
	LogFilename      = "entry_control.log"
)

const Port = ":6061"

const TimeFormat = time.RFC3339

const (
	DatabaseSyncTime      = 5 * time.Minute
	HistoryPrunningTime   = 12 * time.Hour
	HistoryPrunningPeriod = 6 * 24 * time.Hour
)

const (
	UsersBucket   = "Users"
	BookingBucket = "Booking"
)

const (
	DailyCapacity = 5
)

var (
	ErrorUserIdIsNotDefined        = errors.New("User has not any booking")
	ErrorCapacityIsFull            = errors.New("Capacity is full")
	ErrorBookingCalculationFailure = errors.New("Booking calculation failure")
)
