// Features:
// Users should be able to register themselves with a unique identification number.
// Onboard a vaccination center along with mentioned attributes
// Add capacity to a vaccination center per day.
// List All Vaccination Centers with day and capacity details for a given district
// Users should be able to book a center in their district by a day if capacity is available for that day.
// List down all the bookings made for a particular vaccination center.
// Users should be able to cancel the existing booking and vaccination centers should be free to be booked again for that appointment.
// The user should be able to search all the vaccination centers available in the district where the user is currently located. Eg: A user, currently in Karnal, should be able to list down all vaccination centers in Karnal.

package main

import (
	cv "covid-vaccination/services"
	"fmt"
)

type serviceApi interface {
	registerUser(user *cv.User) error
	registerVaccinationCenter(vaccineCenter *cv.Center) error
	addVaccinationCentersCapacity(vaccineCenterId string, day int, capacity int) error
	getVaccinationCenterDetailsOfDistrict(districtId string) ([]cv.Center, error)
	getAllBookingForCenterForDay(centerId string, day int) ([]cv.Booking, error)
	bookCenterinDistrictForTheDay(userId string, centerId string, day int) error
	cancelBookingForUser(userId string) error
}

type serviceO struct {
	userService    *cv.UserService
	bookingService *cv.BookingService
	centerService  *cv.CenterService
}

func getserviceO(userService *cv.UserService, bookingService *cv.BookingService, centerService *cv.CenterService) *serviceO {
	return &serviceO{
		userService:    userService,
		bookingService: bookingService,
		centerService:  centerService,
	}
}

func (s *serviceO) registerUser(user *cv.User) error {
	return s.userService.RegisterUser(user)
}

func (s *serviceO) registerVaccinationCenter(vaccineCenter *cv.Center) error {
	return s.centerService.RegisterCenter(vaccineCenter)
}

func (s *serviceO) addVaccinationCentersCapacity(vaccineCenterId string, day int, capacity int) error {
	center, err := s.centerService.GetCenter(vaccineCenterId)
	if err != nil {
		return err
	}
	center.BookingCapacityPerDay[day] = capacity
	return nil
}

func (s *serviceO) getVaccinationCenterDetailsOfDistrict(districtId string) ([]cv.Center, error) {
	centers, err := s.centerService.GetCentersFromDistrict(districtId)
	if err != nil {
		return nil, err
	}
	return centers, nil
}

func (s *serviceO) getAllBookingForCenterForDay(centerId string, day int) ([]cv.Booking, error) {
	bookings, err := s.bookingService.GetBookingsForCenterForDay(centerId, day)
	if err != nil {
		return nil, err
	}
	return bookings, nil
}

func (s *serviceO) bookCenterinDistrictForTheDay(userId string, centerId string, day int) error {
	return s.bookingService.RegisterBooking(userId, centerId, day)
}

func (s *serviceO) cancelBookingForUser(userId string) error {
	return s.bookingService.CancelBooking(userId)
}

// drive main function
func main() {

	userService := cv.GetUserService(10000)
	centerService := cv.GetCenterService(10000)
	bookingService := cv.GetBookingService(10000, userService, centerService)

	s := getserviceO(userService, bookingService, centerService)

	s.registerUser(&cv.User{
		Id:         "U1",
		Name:       "Harry",
		Age:        35,
		DistrictId: "d1",
	})

	s.registerUser(&cv.User{
		Id:         "U2",
		Name:       "Ron",
		Age:        30,
		DistrictId: "d1",
	})

	s.registerUser(&cv.User{
		Id:         "U3",
		Name:       "Albus",
		Age:        30,
		DistrictId: "d1",
	})

	s.registerUser(&cv.User{
		Id:         "U4",
		Name:       "Draco",
		Age:        15,
		DistrictId: "d1",
	})

	s.registerUser(&cv.User{
		Id:         "U5",
		Name:       "Dobby",
		Age:        30,
		DistrictId: "d2",
	})

	s.registerVaccinationCenter(&cv.Center{
		Id:                    "vc1",
		DefaultCapacity:       5,
		DistrictId:            "d1",
		BookingCapacityPerDay: make(map[int]int),
		BookingCountPerDay:    map[int]int{},
	})

	s.registerVaccinationCenter(&cv.Center{
		Id:                    "vc2",
		DefaultCapacity:       5,
		DistrictId:            "d1",
		BookingCapacityPerDay: make(map[int]int),
		BookingCountPerDay:    map[int]int{},
	})

	s.registerVaccinationCenter(&cv.Center{
		Id:                    "vc3",
		DefaultCapacity:       5,
		DistrictId:            "d3",
		BookingCapacityPerDay: make(map[int]int),
		BookingCountPerDay:    map[int]int{},
	})
	s.addVaccinationCentersCapacity("vc1", 1, 1)
	s.addVaccinationCentersCapacity("vc2", 1, 3)
	s.addVaccinationCentersCapacity("vc1", 5, 10)
	s.addVaccinationCentersCapacity("vc3", 3, 4)

	fmt.Println(s.getVaccinationCenterDetailsOfDistrict("d1"))
	s.bookCenterinDistrictForTheDay("U1", "vc1", 1)
	fmt.Println(s.getAllBookingForCenterForDay("vc1", 1))
	s.bookCenterinDistrictForTheDay("U2", "vc2", 1)
	s.bookCenterinDistrictForTheDay("U3", "vc2", 1)
	fmt.Println(s.getAllBookingForCenterForDay("vc2", 1))
	s.bookCenterinDistrictForTheDay("U5", "vc1", 1)
}

// ADD_USER U1 Harry Male 35 Karnataka Bangalore
// ADD_USER U2 Ron Male 30 Karnataka Bangalore
// ADD_USER U3 Albus Male 30 Karnataka Bangalore
// ADD_USER U4 Draco Male 15 Karnataka Bangalore
// ADD_USER U5 Dobby Male 30 Gujarat Surat
// ADD_VACCINATION_CENTER Karnataka Bangalore VC1
// ADD_VACCINATION_CENTER Karnataka Bangalore VC2
// ADD_VACCINATION_CENTER Karnataka Mysore VC3
// ADD_CAPACITY VC1 1 1
// ADD_CAPACITY VC2 1 3
// ADD_CAPACITY VC1 5 10
// ADD_CAPACITY VC3 3 4
// LIST_VACCINATION_CENTERS Bangalore
// VC1 1 1
// VC1 5 10
// VC2 1 3
// BOOK_VACCINATION VC1 1 U1
// LIST_ALL_BOOKINGS VC1 1
//  BK1 Harry VC1 Bangalore
// BOOK_VACCINATION VC2 1 U2
// BOOK_VACCINATION VC2 1 U3
// LIST_ALL_BOOKINGS VC2 1
//  BK2 Ron VC2 Bangalore
//  BK3 Albus VC2 Bangalore
// BOOK_VACCINATION VC1 1 U5
