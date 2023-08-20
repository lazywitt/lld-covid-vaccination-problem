package covidvaccination

import (
	"fmt"

	"github.com/google/uuid"
)

type BookingService struct {
	BookingData          []*Booking
	BookingCount         int
	BookingIdToIndexMap  map[string]int
	UserIdToBookingIdMap map[string]string
	UserService          *UserService
	CenterService        *CenterService
}

func GetBookingService(bookingSize int, userService *UserService, centerService *CenterService) *BookingService {
	return &BookingService{
		BookingData:          make([]*Booking, bookingSize),
		BookingIdToIndexMap:  make(map[string]int),
		UserIdToBookingIdMap: make(map[string]string),
		UserService: userService,
		CenterService: centerService,
	}
}

type Booking struct {
	Id          string
	Name        string
	CenterId    string
	UserId      string
	Day         int
	IsCancelled bool
}

func (s *BookingService) GetBooking(userId string) (*Booking, error) {
	if _, ok := s.UserIdToBookingIdMap[userId]; !ok {
		return nil, fmt.Errorf("user has no booking")
	}
	bookingId := s.UserIdToBookingIdMap[userId]
	if _, ok := s.BookingIdToIndexMap[bookingId]; !ok {
		return nil, fmt.Errorf("booking Id has no entry mapped in Data")
	}
	return s.BookingData[s.BookingIdToIndexMap[bookingId]], nil
}

func (s *BookingService) RegisterBooking(userId string, centerId string, day int) error {
	user, err := s.UserService.GetUser(userId)
	if err != nil {
		return err
	}
	if user.Age < 18 {
		return fmt.Errorf("user is underage for vaccine")
	}
	if user.HasBooking {
		return fmt.Errorf("user already has a booking")
	}
	center, err := s.CenterService.GetCenter(centerId)
	if err != nil {
		return err
	}
	capacity := center.DefaultCapacity
	if _, ok := center.BookingCapacityPerDay[day]; ok {
		capacity = center.BookingCapacityPerDay[day]
	}
	bookingCount := 0
	if _, ok := center.BookingCountPerDay[day]; ok {
		bookingCount = center.BookingCountPerDay[day]
	}

	if capacity > bookingCount {
		center.BookingCountPerDay[day]++
		user.HasBooking = true
		newBooking := &Booking{
			Id:       uuid.New().String(),
			Name:     user.Name,
			CenterId: centerId,
			UserId:   userId,
			Day:      day,
		}
		s.BookingData[s.BookingCount] = newBooking
		s.BookingIdToIndexMap[newBooking.Id] = s.BookingCount
		s.UserIdToBookingIdMap[userId] = newBooking.Id
		s.BookingCount++
	} else {
		return fmt.Errorf("no empty slot left for booking")
	}

	return nil
}

func (s *BookingService) CancelBooking(userId string) error {
	user, err := s.UserService.GetUser(userId)
	if err != nil {
		return err
	}
	if !user.HasBooking {
		return fmt.Errorf("user has no booking")
	}
	booking, err := s.GetBooking(userId)
	if err != nil {
		return err
	}
	user.HasBooking = false
	booking.IsCancelled = true
	return nil
}

func (s *BookingService) GetBookingsForCenterForDay(centerId string, day int) ([]Booking, error) {
	bookings := []Booking{}
	for i := 0; i < s.BookingCount; i++ {
		booking := s.BookingData[i]
		if booking.CenterId == centerId && booking.Day == day && !booking.IsCancelled {
			bookings = append(bookings, *booking)
		}
	}
	return bookings, nil
}
