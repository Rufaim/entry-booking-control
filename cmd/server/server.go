package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/Rufaim/entry_booking_control/cmd/message"
	"github.com/boltdb/bolt"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)

type Server struct {
	pb.UnimplementedLabVisitsServiceServer
	DB     *bolt.DB
	logger *log.Logger
}

func (s *Server) SetUserVisit(ctx context.Context, uvs *pb.UserVisitSet) (*pb.UserVisitSetResult, error) {
	s.logger.Printf("Updating book on %s for %s\n", uvs.GetDay(), uvs.User.GetId())
	userID := []byte(uvs.User.Id)
	err := s.DB.Update(func(tx *bolt.Tx) error {
		bookBucket := tx.Bucket([]byte(BookingBucket))
		booked, err := bytesToInt(bookBucket.Get([]byte(uvs.Day.String())))
		if err != nil {
			return err
		}
		if booked >= DailyCapacity {
			return ErrorCapacityIsFull
		}
		userBucket := tx.Bucket([]byte(UsersBucket))
		userHistoryByte := userBucket.Get(userID)
		userHistory := &pb.UserStat{}
		if userHistoryByte == nil {
			userHistory = &pb.UserStat{
				User: uvs.User,
			}
		} else {
			if err := proto.Unmarshal(userHistoryByte, userHistory); err != nil {
				return err
			}
		}
		bookingUpdated := false
		for _, visit := range userHistory.Visits {
			if visit.Day == uvs.Day {
				visit.Timestamp = uvs.Timestamp
				bookingUpdated = true
				break
			}
		}
		if !bookingUpdated {
			userHistory.Visits = append(userHistory.Visits, &pb.Visit{
				Day:       uvs.Day,
				Timestamp: uvs.Timestamp,
			})
			if err := bookBucket.Put([]byte(uvs.Day.String()), intToBytes(booked+1)); err != nil {
				return err
			}
		}
		message, err := proto.Marshal(userHistory)
		if err != nil {
			return err
		}

		return userBucket.Put(userID, message)
	})

	if err == ErrorCapacityIsFull {
		return &pb.UserVisitSetResult{
			Status: pb.UserVisitSetResult_FAILURE,
			Text:   fmt.Sprintf("Lab is fully booked for %s", uvs.Day.String()),
		}, nil
	}

	if err != nil {
		return &pb.UserVisitSetResult{
			Status: pb.UserVisitSetResult_FAILURE,
			Text:   err.Error(),
		}, err
	}
	return &pb.UserVisitSetResult{
		Status: pb.UserVisitSetResult_OK,
	}, nil
}

func (s *Server) DelUserVisit(ctx context.Context, uvs *pb.UserVisitSet) (*pb.UserVisitSetResult, error) {
	s.logger.Printf("Deleting book on %s for %s\n", uvs.GetDay(), uvs.User.GetId())
	userID := []byte(uvs.User.Id)
	err := s.DB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(UsersBucket))
		userHistoryByte := bucket.Get(userID)
		if userHistoryByte == nil {
			return ErrorUserIdIsNotDefined
		}
		userHistory := &pb.UserStat{}
		err := proto.Unmarshal(userHistoryByte, userHistory)
		if err != nil {
			return err
		}

		isDeleted := false
		for i, visit := range userHistory.Visits {
			if visit.Day == uvs.Day {
				// this is how deleting of element works in go
				userHistory.Visits[i] = userHistory.Visits[len(userHistory.Visits)-1]
				userHistory.Visits = userHistory.Visits[:len(userHistory.Visits)-1]
				isDeleted = true
				break
			}
		}

		if isDeleted {
			message, err := proto.Marshal(userHistory)
			if err != nil {
				return err
			}
			if err := bucket.Put(userID, message); err != nil {
				return err
			}

			bucket = tx.Bucket([]byte(BookingBucket))
			booked, err := bytesToInt(bucket.Get([]byte(uvs.Day.String())))
			fmt.Println("Booked ", booked)
			if err != nil {
				return err
			}
			if booked <= 0 {
				return ErrorBookingCalculationFailure
			}
			return bucket.Put([]byte(uvs.Day.String()), intToBytes(booked-1))
		}
		return nil
	})

	if err == ErrorUserIdIsNotDefined {
		return &pb.UserVisitSetResult{
			Status: pb.UserVisitSetResult_FAILURE,
			Text:   err.Error(),
		}, nil
	}

	if err != nil {
		return &pb.UserVisitSetResult{
			Status: pb.UserVisitSetResult_FAILURE,
			Text:   err.Error(),
		}, err
	}

	return &pb.UserVisitSetResult{
		Status: pb.UserVisitSetResult_OK,
	}, nil
}

func (s *Server) GetVisits(ctx context.Context, usr *pb.User) (*pb.UserStat, error) {
	s.logger.Printf("Report on %s requested\n", usr.GetId())
	userID := []byte(usr.GetId())
	var visits []*pb.Visit
	err := s.DB.View(func(tx *bolt.Tx) error {
		bucker := tx.Bucket([]byte(UsersBucket))
		userHistoryByte := bucker.Get(userID)
		if userHistoryByte == nil {
			return ErrorUserIdIsNotDefined
		}
		userHistory := &pb.UserStat{}
		if err := proto.Unmarshal(userHistoryByte, userHistory); err != nil {
			return err
		}
		visits = userHistory.Visits
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &pb.UserStat{
		User:   usr,
		Visits: visits,
	}, nil
}

func (s *Server) GetAllVisitsReport(ctx context.Context, in *pb.Void) (*pb.AllUsersStat, error) {
	s.logger.Println("Full report requested")
	stat := make([]*pb.UserStat, 0)
	books := make(map[string]uint32)
	err := s.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(UsersBucket))
		c := bucket.Cursor()
		for userID, userHistoryByte := c.First(); userID != nil; userID, userHistoryByte = c.Next() {
			userHistory := &pb.UserStat{}
			if err := proto.Unmarshal(userHistoryByte, userHistory); err != nil {
				return err
			}
			stat = append(stat, userHistory)
		}

		bucket = tx.Bucket([]byte(BookingBucket))
		for name, _ := range pb.Weekday_value {
			t, err := bytesToInt(bucket.Get([]byte(name)))
			if err != nil {
				return err
			}
			books[name] = uint32(t)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &pb.AllUsersStat{
		Stat:          stat,
		BookingAmount: books,
	}, nil
}

func (s *Server) Sync() error {
	s.logger.Println("Server syncronization")
	return s.DB.Sync()
}

func (s *Server) PruneHistory() error {
	s.logger.Println("Prunning history")
	return s.DB.Update(func(tx *bolt.Tx) error {
		currentTime := time.Now()
		bucket := tx.Bucket([]byte(UsersBucket))

		bookingNew := make([]int, len(pb.Weekday_name))
		c := bucket.Cursor()
		for userID, userHistoryByte := c.First(); userID != nil; userID, userHistoryByte = c.Next() {
			userHistory := &pb.UserStat{}
			if err := proto.Unmarshal(userHistoryByte, userHistory); err != nil {
				return err
			}
			newHistory := make([]*pb.Visit, 0)
			for _, visit := range userHistory.Visits {
				timestamp, err := time.Parse(TimeFormat, visit.Timestamp)
				if err != nil {
					return err
				}
				if currentTime.Sub(timestamp) < HistoryPrunningPeriod {
					newHistory = append(newHistory, visit)
				}
			}
			if len(newHistory) == 0 {
				if err := c.Delete(); err != nil {
					return err
				}
				continue
			}

			for _, visit := range newHistory {
				bookingNew[int(visit.Day)]++
			}

			userHistory.Visits = newHistory
			message, err := proto.Marshal(userHistory)
			if err != nil {
				return err
			}
			if err := bucket.Put(userID, message); err != nil {
				return err
			}
		}

		bucket = tx.Bucket([]byte(BookingBucket))
		for name, val := range pb.Weekday_value {
			if err := bucket.Put([]byte(name), intToBytes(bookingNew[val])); err != nil {
				return err
			}
		}

		return nil
	})
}

func NewServer(logger *log.Logger) (*Server, error) {
	logger.Println("Initializing server")
	db, err := bolt.Open(DatabaseFilename, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, errors.Wrap(err, "Database opening failure")
	}
	logger.Println("Database file opened")
	err = db.Update(func(tx *bolt.Tx) error {
		for _, b := range []string{UsersBucket, BookingBucket} {
			_, err := tx.CreateBucketIfNotExists([]byte(b))
			if err != nil {
				return err
			}
		}

		bucket := tx.Bucket([]byte(BookingBucket))
		for _, weekday := range pb.Weekday_name {
			val := bucket.Get([]byte(weekday))
			if val == nil {
				if err := bucket.Put([]byte(weekday), intToBytes(0)); err != nil {
					return err
				}
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	logger.Println("Database initialized")
	return &Server{
		DB:     db,
		logger: logger,
	}, nil
}
